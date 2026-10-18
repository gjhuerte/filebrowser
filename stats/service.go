package stats

import (
	"context"
	"sync"
	"time"
)

// Status is what the API reports about a user's scan.
type Status struct {
	// Scanning is true while a scan is queued or running.
	Scanning bool `json:"scanning"`
	// Summary is the last completed result, or nil if there is none yet.
	Summary *Summary `json:"summary"`
	// ScannedAt is when Summary was completed.
	ScannedAt *time.Time `json:"scannedAt,omitempty"`
	// DurationMs is how long that scan took.
	DurationMs int64 `json:"durationMs,omitempty"`
	// Progress counts what the running scan has found so far.
	Progress ProgressCounts `json:"progress"`
	// Error describes why the last scan failed, if it did.
	Error string `json:"error,omitempty"`
}

// ProgressCounts is a snapshot of Progress.
type ProgressCounts struct {
	Files   int64 `json:"files"`
	Folders int64 `json:"folders"`
}

// Service runs scans in the background and remembers each key's last result,
// so that clients poll for a result instead of holding a request open for the
// whole scan (which reverse proxies time out on large trees).
//
// At most one scan runs per key, and at most maxConcurrent across all keys;
// the rest wait their turn. It is safe for concurrent use.
type Service struct {
	mu      sync.Mutex
	jobs    map[string]*job
	slots   chan struct{}
	timeout time.Duration
}

type job struct {
	running   bool
	progress  *Progress
	summary   *Summary
	scannedAt time.Time
	duration  time.Duration
	err       string
}

// NewService returns a Service running at most maxConcurrent scans at once,
// each stopped after timeout.
func NewService(maxConcurrent int, timeout time.Duration) *Service {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	return &Service{
		jobs:    map[string]*job{},
		slots:   make(chan struct{}, maxConcurrent),
		timeout: timeout,
	}
}

// Status reports the state of key's scan.
func (s *Service) Status(key string) Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.statusLocked(s.jobs[key])
}

// Start begins a scan for key with opts unless one is already running, and
// returns the resulting status. The scan is not tied to any request: it keeps
// going if the client goes away.
func (s *Service) Start(key string, opts Options) Status {
	s.mu.Lock()
	defer s.mu.Unlock()

	j := s.jobs[key]
	if j == nil {
		j = &job{}
		s.jobs[key] = j
	}
	if !j.running {
		j.running = true
		j.progress = &Progress{}
		go s.run(j, opts, j.progress)
	}
	return s.statusLocked(j)
}

func (s *Service) run(j *job, opts Options, progress *Progress) {
	s.slots <- struct{}{}
	defer func() { <-s.slots }()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	started := time.Now()
	summary, err := Scan(ctx, opts, progress)

	s.mu.Lock()
	defer s.mu.Unlock()
	j.running = false
	if err != nil {
		// Keep showing the previous result; just report the failure.
		j.err = err.Error()
		return
	}
	j.summary = summary
	j.scannedAt = time.Now()
	j.duration = j.scannedAt.Sub(started)
	j.err = ""
}

func (s *Service) statusLocked(j *job) Status {
	if j == nil {
		return Status{}
	}
	st := Status{
		Scanning: j.running,
		Summary:  j.summary,
		Error:    j.err,
	}
	if j.summary != nil {
		at := j.scannedAt
		st.ScannedAt = &at
		st.DurationMs = j.duration.Milliseconds()
	}
	if j.running && j.progress != nil {
		st.Progress = ProgressCounts{
			Files:   j.progress.Files.Load(),
			Folders: j.progress.Folders.Load(),
		}
	}
	return st
}
