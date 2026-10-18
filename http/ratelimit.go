package fbhttp

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tomasen/realip"
)

// Failed password attempts allowed per client within failureWindow before
// further attempts are refused with 429 Too Many Requests.
const (
	maxFailedAttempts = 10
	failureWindow     = 15 * time.Minute

	// maxTrackedClients bounds the limiter's memory. When it is reached, entries
	// whose window has passed are dropped first.
	maxTrackedClients = 10000
)

// attemptLimiter counts failed password attempts per key (a client address,
// optionally combined with what was being unlocked) and blocks a key once it
// has failed too often within the window. Successful attempts clear the count.
//
// It is in-memory and per process: restarting File Browser resets it. That is
// enough to stop online guessing; it is not meant as an audit log.
type attemptLimiter struct {
	mu      sync.Mutex
	entries map[string]*attempts
	max     int
	window  time.Duration
	now     func() time.Time
}

type attempts struct {
	count int
	first time.Time
}

func newAttemptLimiter(maxFailures int, window time.Duration) *attemptLimiter {
	return &attemptLimiter{
		entries: map[string]*attempts{},
		max:     maxFailures,
		window:  window,
		now:     time.Now,
	}
}

// loginLimiter guards /api/login; shareLimiter guards password-protected shares.
var (
	loginLimiter = newAttemptLimiter(maxFailedAttempts, failureWindow)
	shareLimiter = newAttemptLimiter(maxFailedAttempts, failureWindow)
)

// Allow reports whether key may try again and, if not, how long until it can.
func (l *attemptLimiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		return true, 0
	}

	elapsed := l.now().Sub(e.first)
	if elapsed >= l.window {
		delete(l.entries, key)
		return true, 0
	}

	if e.count >= l.max {
		return false, l.window - elapsed
	}

	return true, 0
}

// Fail records a failed attempt for key.
func (l *attemptLimiter) Fail(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	e, ok := l.entries[key]
	if !ok || now.Sub(e.first) >= l.window {
		if !ok && len(l.entries) >= maxTrackedClients {
			l.pruneLocked(now)
		}
		l.entries[key] = &attempts{count: 1, first: now}
		return
	}

	e.count++
}

// Reset forgets the failures recorded for key, e.g. after a successful login.
func (l *attemptLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}

// pruneLocked drops expired entries, and if that is not enough to make room,
// all of them. The caller must hold l.mu.
func (l *attemptLimiter) pruneLocked(now time.Time) {
	for k, e := range l.entries {
		if now.Sub(e.first) >= l.window {
			delete(l.entries, k)
		}
	}
	if len(l.entries) >= maxTrackedClients {
		l.entries = map[string]*attempts{}
	}
}

// clientKey identifies the client for rate limiting.
//
// realip honours X-Real-IP and X-Forwarded-For, which is what makes the limit
// per visitor when File Browser runs behind a reverse proxy. When it is exposed
// directly, a client can forge those headers to spread its attempts over many
// keys, so put File Browser behind a proxy that sets them.
func clientKey(r *http.Request) string {
	return realip.FromRequest(r)
}

// tooManyAttempts writes the Retry-After header and returns the status to send
// when a client has been blocked.
func tooManyAttempts(w http.ResponseWriter, retryAfter time.Duration) int {
	seconds := int(retryAfter.Round(time.Second) / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	w.Header().Set("Retry-After", strconv.Itoa(seconds))
	return http.StatusTooManyRequests
}
