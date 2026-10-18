package fbhttp

import (
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/gorilla/websocket"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/storage/bolt"
	"github.com/filebrowser/filebrowser/v2/users"
)

// A command that writes a lot to stderr before printing to stdout used to hang
// the terminal: stdout was read to the end first while the command blocked on
// a full stderr pipe.
func TestCommandsHandlerStreamsLargeStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/sh")
	}

	db, err := storm.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	st, err := bolt.NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}
	key := []byte("test-signing-key")
	if err := st.Settings.Save(&settings.Settings{Key: key}); err != nil {
		t.Fatal(err)
	}
	perm := users.Permissions{Execute: true}
	if err := st.Users.Save(&users.User{Username: "u", Password: "pw", Perm: perm, Commands: []string{"sh"}}); err != nil {
		t.Fatal(err)
	}

	server := &settings.Server{Root: t.TempDir(), EnableExec: true}
	srv := httptest.NewServer(handle(commandsHandler, "", st, server))
	t.Cleanup(srv.Close)

	header := map[string][]string{"X-Auth": {signToken(t, perm, key)}}
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/", header)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	cmd := `sh -c "head -c 200000 /dev/zero | tr '\\000' x >&2; echo done"`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
		t.Fatal(err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("terminal output ended before \"done\" (hung or failed): %v", err)
		}
		if string(msg) == "done" {
			return
		}
	}
}
