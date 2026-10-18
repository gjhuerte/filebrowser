package fbhttp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/filebrowser/filebrowser/v2/runner"
)

const (
	// WSWriteDeadline bounds every write to the terminal WebSocket, so a client
	// that stops reading cannot stall the command forever.
	WSWriteDeadline = 10 * time.Second

	// maxOutputLine is the longest single line of command output sent as one
	// message. Longer lines end the live output; the rest is discarded so the
	// command can still finish.
	maxOutputLine = 1 << 20 // 1 MiB

	// commandWaitDelay is how long a stopped command's children may keep its
	// output open before File Browser stops waiting for them.
	commandWaitDelay = 2 * time.Second
)

// upgrader keeps gorilla/websocket's default origin check, which rejects
// cross-site WebSocket requests that would otherwise ride on the auth cookie.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	cmdNotAllowed = []byte("Command not allowed.")
)

// wsErr logs a terminal failure and closes the WebSocket with a short reason.
func wsErr(ws *websocket.Conn, r *http.Request, status int, err error) {
	txt := http.StatusText(status)
	if err != nil || status >= 400 {
		log.Printf("%s: %v %s %v", r.URL.Path, status, r.RemoteAddr, err)
	}
	if err := ws.WriteControl(websocket.CloseInternalServerErr, []byte(txt), time.Now().Add(WSWriteDeadline)); err != nil {
		log.Print(err)
	}
}

// wsWrite sends one text message, bounded by WSWriteDeadline.
func wsWrite(ws *websocket.Conn, msg []byte) error {
	if err := ws.SetWriteDeadline(time.Now().Add(WSWriteDeadline)); err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, msg)
}

// commandsHandler runs one allowed command per WebSocket connection in the
// folder named by the request path and streams its output back line by line.
//
// The client sends the command line as the first message. Closing the socket
// stops the command: the process is killed so nothing is left running after
// the user leaves the page or cancels.
var commandsHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	var raw string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
			return 0, nil
		}

		raw = strings.TrimSpace(string(msg))
		if raw != "" {
			break
		}
	}

	// Fail fast
	if !d.server.EnableExec || !d.user.Perm.Execute || !d.Check(r.URL.Path) {
		if err := wsWrite(conn, cmdNotAllowed); err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
		}

		return 0, nil
	}

	command, name, err := runner.ParseCommand(d.settings, raw)
	if err != nil {
		if err := wsWrite(conn, []byte(err.Error())); err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
		}
		return 0, nil
	}

	if !slices.Contains(d.user.Commands, name) {
		if err := wsWrite(conn, cmdNotAllowed); err != nil {
			wsErr(conn, r, http.StatusInternalServerError, err)
		}

		return 0, nil
	}

	// The request context is not cancelled when a hijacked WebSocket closes,
	// so watch the socket instead: any read error means the client is gone.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				cancel()
				return
			}
		}
	}()

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Dir = d.user.FullPath(r.URL.Path)
	cmd.WaitDelay = commandWaitDelay

	// One pipe for both streams, so a command writing a lot to stderr cannot
	// block while stdout is still being read (and the other way round).
	out, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
		return 0, nil
	}

	waitErr := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		_ = pw.Close()
		waitErr <- err
	}()

	s := bufio.NewScanner(out)
	s.Buffer(make([]byte, 0, 64*1024), maxOutputLine)
	for s.Scan() {
		if err := wsWrite(conn, s.Bytes()); err != nil {
			// The client is gone; stop the command and just drain its output.
			cancel()
			break
		}
	}
	// Keep reading until the command exits so it never blocks on a full pipe.
	_, _ = io.Copy(io.Discard, out)

	if err := <-waitErr; err != nil && ctx.Err() == nil {
		wsErr(conn, r, http.StatusInternalServerError, err)
	}

	return 0, nil
})
