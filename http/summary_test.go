package fbhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/stats"
	"github.com/filebrowser/filebrowser/v2/users"
)

// The Home summary is started with POST and polled with GET, so the scan
// keeps going after the starting request has returned.
func TestSummaryStartAndPoll(t *testing.T) {
	userScope := t.TempDir()
	if err := os.MkdirAll(filepath.Join(userScope, "Photos"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"Photos/a.jpg": "jpeg", "notes.txt": "hi"} {
		if err := os.WriteFile(filepath.Join(userScope, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// Outside the scope; the link to it must not be followed.
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(userScope, "escape")); err != nil {
		t.Fatal(err)
	}

	key := []byte("test-signing-key")
	perm := users.Permissions{Download: true}
	st := scopedUserStorage(t, userScope, perm, key)
	token := signToken(t, perm, key)
	get := handle(summaryGetHandler, "", st, &settings.Server{})
	post := handle(summaryPostHandler, "", st, &settings.Server{})

	call := func(h http.Handler, method string) stats.Status {
		t.Helper()
		req, _ := http.NewRequest(method, "/api/summary", http.NoBody)
		req.Header.Set("X-Auth", token)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d body = %q", method, rec.Code, rec.Body.String())
		}
		var s stats.Status
		if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return s
	}

	if s := call(post, http.MethodPost); !s.Scanning && s.Summary == nil {
		t.Fatalf("POST neither started a scan nor returned a result: %+v", s)
	}

	deadline := time.Now().Add(5 * time.Second)
	var s stats.Status
	for {
		s = call(get, http.MethodGet)
		if !s.Scanning && s.Summary != nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("scan did not finish: %+v", s)
		}
		time.Sleep(10 * time.Millisecond)
	}
	if s.Summary.Files != 2 || s.Summary.Bytes != 6 {
		t.Fatalf("summary = %+v", s.Summary)
	}
}
