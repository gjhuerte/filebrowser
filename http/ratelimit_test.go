package fbhttp

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/asdine/storm/v3"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
)

func TestAttemptLimiterBlocksAfterMaxFailures(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := newAttemptLimiter(3, time.Minute)
	l.now = func() time.Time { return now }

	for i := 0; i < 3; i++ {
		if ok, _ := l.Allow("ip"); !ok {
			t.Fatalf("attempt %d blocked too early", i+1)
		}
		l.Fail("ip")
	}

	ok, wait := l.Allow("ip")
	if ok {
		t.Fatal("expected the client to be blocked after 3 failures")
	}
	if wait <= 0 || wait > time.Minute {
		t.Errorf("unexpected retry-after %v", wait)
	}

	// Other clients are unaffected.
	if ok, _ := l.Allow("other"); !ok {
		t.Error("an unrelated client was blocked")
	}

	// The block lifts once the window has passed.
	now = now.Add(time.Minute)
	if ok, _ := l.Allow("ip"); !ok {
		t.Error("client still blocked after the window passed")
	}
}

func TestAttemptLimiterResetClearsFailures(t *testing.T) {
	l := newAttemptLimiter(2, time.Minute)
	l.Fail("ip")
	l.Fail("ip")
	l.Reset("ip")
	if ok, _ := l.Allow("ip"); !ok {
		t.Error("client still blocked after Reset")
	}
}

func TestAttemptLimiterBoundsMemory(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	l := newAttemptLimiter(1, time.Minute)
	l.now = func() time.Time { return now }

	for i := 0; i < maxTrackedClients; i++ {
		l.Fail("client-" + strconv.Itoa(i))
	}
	now = now.Add(2 * time.Minute)
	l.Fail("new")

	if len(l.entries) > maxTrackedClients {
		t.Errorf("limiter grew past its bound: %d entries", len(l.entries))
	}
}

// After too many wrong passwords, /api/login answers 429 with Retry-After,
// even for the right password, until the window passes.
func TestLoginHandlerRateLimitsFailedAttempts(t *testing.T) {
	old := loginLimiter
	loginLimiter = newAttemptLimiter(maxFailedAttempts, failureWindow)
	t.Cleanup(func() { loginLimiter = old })

	db, err := storm.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	st, err := bolt.NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Settings.Save(&settings.Settings{Key: []byte("key"), AuthMethod: auth.MethodJSONAuth}); err != nil {
		t.Fatal(err)
	}
	if err := st.Auth.Save(&auth.JSONAuth{}); err != nil {
		t.Fatal(err)
	}
	hash, err := users.HashPwd("correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Users.Save(&users.User{Username: "admin", Password: hash}); err != nil {
		t.Fatal(err)
	}

	server := &settings.Server{Root: t.TempDir()}
	login := func(password string) *httptest.ResponseRecorder {
		body := `{"username":"admin","password":"` + password + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
		req.RemoteAddr = "203.0.113.7:4321"
		rec := httptest.NewRecorder()
		handle(loginHandler(time.Hour), "", st, server).ServeHTTP(rec, req)
		return rec
	}

	for i := 0; i < maxFailedAttempts; i++ {
		if rec := login("wrong"); rec.Code != http.StatusForbidden {
			t.Fatalf("attempt %d: expected 403, got %d", i+1, rec.Code)
		}
	}

	rec := login("correct-horse")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after %d failures, got %d", maxFailedAttempts, rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("429 response is missing Retry-After")
	}
}
