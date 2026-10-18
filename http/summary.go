package fbhttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/filebrowser/filebrowser/v2/stats"
)

// summaryService runs the Home dashboard scans. Two scans at most run at once
// across all users, so a burst of refreshes cannot saturate the disk, and a
// scan that takes over an hour is abandoned.
var summaryService = stats.NewService(2, time.Hour)

// summaryKey identifies a user's cached result. The scope is part of it, so
// changing a user's scope does not show the old folder's numbers.
func summaryKey(d *data) string {
	return strconv.FormatUint(uint64(d.user.ID), 10) + ":" + d.user.FullPath("/")
}

// summaryGetHandler returns the state of the user's Home dashboard scan: the
// last result, whether a scan is running and how far it has got.
var summaryGetHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return renderJSON(w, r, summaryService.Status(summaryKey(d)))
})

// summaryPostHandler starts a scan of the user's whole scope, unless one is
// already running, and returns at once. Clients poll summaryGetHandler for
// the result, so a long scan never hits a reverse proxy's request timeout.
var summaryPostHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// Copy what the scan needs: it outlives this request, and d with it.
	user := *d.user
	check := &data{settings: d.settings, server: d.server, user: &user}

	status := summaryService.Start(summaryKey(d), stats.Options{
		Root: user.FullPath("/"),
		// Rules and hidden dotfiles apply, as in the file list.
		Allow: check.Check,
	})
	return renderJSON(w, r, status)
})
