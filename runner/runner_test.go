package runner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/users"
)

// A hook run through a shell must treat the file name as data. Before the fix,
// $FILE was pasted into the script text, so a name containing ";" ran a second
// command.
func TestRunHookShellDoesNotInjectFileName(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/sh")
	}

	scope := t.TempDir()
	out := filepath.Join(t.TempDir(), "out.txt")
	marker := filepath.Join(t.TempDir(), "pwned")

	r := &Runner{
		Enabled: true,
		Settings: &settings.Settings{
			Shell: []string{"sh", "-c"},
			Commands: map[string][]string{
				"after_upload": {`printf '%s' "$FILE" > ` + out},
			},
		},
	}

	user := &users.User{Username: "u", Password: "x", Scope: "/"}
	if err := user.Clean(scope, false); err != nil {
		t.Fatal(err)
	}

	name := "/a;touch " + marker
	if err := r.RunHook(func() error { return nil }, "upload", name, "", user); err != nil {
		t.Fatalf("hook failed: %v", err)
	}

	if _, err := os.Stat(marker); err == nil {
		t.Fatal("VULNERABLE: the file name was executed as a shell command")
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("hook did not run: %v", err)
	}
	if !strings.HasSuffix(string(got), name) {
		t.Errorf("hook saw FILE=%q, want it to end with %q", got, name)
	}
}

// Without a shell, variables are still expanded into the program's arguments.
func TestRunHookWithoutShellExpandsArguments(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses touch")
	}

	scope := t.TempDir()
	r := &Runner{
		Enabled: true,
		Settings: &settings.Settings{
			Commands: map[string][]string{
				"after_upload": {"touch $FILE.done"},
			},
		},
	}

	user := &users.User{Username: "u", Password: "x", Scope: "/"}
	if err := user.Clean(scope, false); err != nil {
		t.Fatal(err)
	}

	if err := r.RunHook(func() error { return nil }, "upload", "/report", "", user); err != nil {
		t.Fatalf("hook failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(scope, "report.done")); err != nil {
		t.Errorf("expected $FILE to be expanded into the argument: %v", err)
	}
}
