package stats

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// makeTree creates files (path -> size) under a temp dir and returns it.
func makeTree(t *testing.T, files map[string]int) string {
	t.Helper()
	root := t.TempDir()
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	for p, size := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, make([]byte, size), 0o644); err != nil {
			t.Fatal(err)
		}
		mod := base.Add(time.Duration(size) * time.Minute)
		if err := os.Chtimes(full, mod, mod); err != nil {
			t.Fatal(err)
		}

	}
	return root
}

func TestFileType(t *testing.T) {
	cases := map[string]TypeKey{
		"a.JPG":     TypeImages,
		"clip.flv":  TypeVideos,
		"song.mp3":  TypeAudio,
		"r.pdf":     TypeDocuments,
		"b.tar.gz":  TypeArchives,
		".bashrc":   TypeOther,
		"Makefile":  TypeOther,
		"weird.xyz": TypeOther,
	}
	for name, want := range cases {
		if got := FileType(name); got != want {
			t.Errorf("FileType(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestScanSummary(t *testing.T) {
	root := makeTree(t, map[string]int{
		"Photos/a.jpg":           300,
		"Photos/2025/b.png":      200,
		"Music/song.mp3":         500,
		"notes.txt":              10,
		"Projects/x/y/z/deep.go": 1,
		"Projects/big.zip":       900,
		"Projects/c.bin":         50,
		"secret/hidden.txt":      5000,
	})
	if err := os.Symlink(filepath.Join(root, "Music"), filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}

	progress := &Progress{}
	s, err := Scan(context.Background(), Options{
		Root:    root,
		Allow:   func(v string) bool { return !strings.HasPrefix(v, "/secret") },
		Workers: 3,
	}, progress)
	if err != nil {
		t.Fatal(err)
	}

	if s.Files != 7 {
		t.Errorf("Files = %d, want 7", s.Files)
	}
	// Photos, Photos/2025, Music, Projects, Projects/x, x/y, x/y/z
	if s.Folders != 7 {
		t.Errorf("Folders = %d, want 7", s.Folders)
	}
	if s.Bytes != 1961 {
		t.Errorf("Bytes = %d, want 1961", s.Bytes)
	}
	if progress.Files.Load() != 7 || progress.Folders.Load() != 7 {
		t.Errorf("progress = %d/%d", progress.Files.Load(), progress.Folders.Load())
	}

	if s.ByType[0].Key != TypeArchives || s.ByType[len(s.ByType)-1].Key != TypeOther {
		t.Errorf("ByType order = %+v", s.ByType)
	}
	if s.Largest[0].Path != "/Projects/big.zip" || len(s.Largest) != TopN {
		t.Errorf("Largest = %+v", s.Largest)
	}
	// Modification time grows with size, so the most recent is the biggest.
	if s.Recent[0].Name != "big.zip" {
		t.Errorf("Recent[0] = %+v", s.Recent[0])
	}
	for _, e := range append(s.Recent, s.Largest...) {
		if strings.HasPrefix(e.Path, "/secret") || strings.HasPrefix(e.Path, "/link") {
			t.Errorf("unexpected entry %s", e.Path)
		}
	}
}

func TestScanCancelled(t *testing.T) {
	root := makeTree(t, map[string]int{"a/b.txt": 1})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Scan(ctx, Options{Root: root}, nil); err == nil {
		t.Fatal("expected an error from a cancelled scan")
	}
}

func TestScanManyFolders(t *testing.T) {
	files := map[string]int{}
	for i := 0; i < 40; i++ {
		for j := 0; j < 5; j++ {
			files[filepath.ToSlash(filepath.Join(
				"d"+string(rune('a'+i%26))+string(rune('a'+i/26)),
				"f"+string(rune('a'+j))+".txt"))] = j + 1
		}
	}
	root := makeTree(t, files)
	s, err := Scan(context.Background(), Options{Root: root, Workers: 8}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s.Files != 200 || s.Folders != 40 || s.Bytes != 40*15 {
		t.Errorf("got %d files, %d folders, %d bytes", s.Files, s.Folders, s.Bytes)
	}
}

func TestService(t *testing.T) {
	root := makeTree(t, map[string]int{"a.txt": 3, "b/c.mp3": 4})
	svc := NewService(2, time.Minute)

	if st := svc.Status("u1"); st.Scanning || st.Summary != nil {
		t.Fatalf("fresh status = %+v", st)
	}

	st := svc.Start("u1", Options{Root: root})
	if !st.Scanning {
		t.Fatalf("Start did not report scanning: %+v", st)
	}
	// A second Start while running must not start another scan.
	svc.Start("u1", Options{Root: root})

	deadline := time.Now().Add(5 * time.Second)
	for {
		st = svc.Status("u1")
		if !st.Scanning {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("scan did not finish")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if st.Summary == nil || st.Summary.Files != 2 || st.ScannedAt == nil || st.Error != "" {
		t.Fatalf("finished status = %+v", st)
	}
	if other := svc.Status("u2"); other.Summary != nil {
		t.Fatal("results leaked between keys")
	}

	// A failing scan keeps the previous result.
	short := NewService(1, time.Nanosecond)
	short.jobs["u1"] = &job{summary: st.Summary, scannedAt: *st.ScannedAt}
	short.Start("u1", Options{Root: root})
	for short.Status("u1").Scanning {
		time.Sleep(time.Millisecond)
	}
	if got := short.Status("u1"); got.Summary == nil || got.Error == "" {
		t.Fatalf("after failure = %+v", got)
	}
}
