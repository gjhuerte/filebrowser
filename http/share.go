package fbhttp

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/users"
	"golang.org/x/crypto/bcrypt"
)

// shareResponse is the client-facing representation of a share. It deliberately
// omits the server-side secrets of share.Link — the bcrypt PasswordHash (which
// would be crackable offline) and the bypass Token — exposing only whether the
// share is password-protected via HasPassword.
type shareResponse struct {
	Hash        string `json:"hash"`
	Path        string `json:"path"`
	UserID      uint   `json:"userID"`
	Expire      int64  `json:"expire"`
	HasPassword bool   `json:"hasPassword"`
}

// toShareResponse converts a link to its client-facing form.
func toShareResponse(l *share.Link) *shareResponse {
	return &shareResponse{
		Hash:        l.Hash,
		Path:        l.Path,
		UserID:      l.UserID,
		Expire:      l.Expire,
		HasPassword: l.PasswordHash != "",
	}
}

// toShareResponses converts links to their client-facing form.
func toShareResponses(links []*share.Link) []*shareResponse {
	res := make([]*shareResponse, 0, len(links))
	for _, l := range links {
		res = append(res, toShareResponse(l))
	}
	return res
}

// withPermShare requires the Share and Download permissions.
func withPermShare(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if !d.user.Perm.Share || !d.user.Perm.Download {
			return http.StatusForbidden, nil
		}

		return fn(w, r, d)
	})
}

// shareListHandler serves GET /api/shares: every link for admins, otherwise
// the caller's own.
var shareListHandler = withPermShare(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var (
		s   []*share.Link
		err error
	)
	if d.user.Perm.Admin {
		s, err = d.store.Share.All()
	} else {
		s, err = d.store.Share.FindByUserID(d.user.ID)
	}
	if errors.Is(err, fberrors.ErrNotExist) {
		return renderJSON(w, r, []*shareResponse{})
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	sort.Slice(s, func(i, j int) bool {
		if s[i].UserID != s[j].UserID {
			return s[i].UserID < s[j].UserID
		}
		return s[i].Expire < s[j].Expire
	})

	return renderJSON(w, r, toShareResponses(s))
})

// shareGetsHandler serves GET /api/share/{path}: the links for one file.
var shareGetsHandler = withPermShare(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	var (
		s   []*share.Link
		err error
	)
	if d.user.Perm.Admin {
		s, err = getSharesForAdminPath(d, r.URL.Path)
	} else {
		s, err = d.store.Share.Gets(r.URL.Path, d.user.ID)
	}
	if errors.Is(err, fberrors.ErrNotExist) {
		return renderJSON(w, r, []*shareResponse{})
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, toShareResponses(s))
})

// getSharesForAdminPath finds the links, by any user, that point at the file
// the admin sees at path. Owners have different scopes, so paths are compared
// on disk.
func getSharesForAdminPath(d *data, path string) ([]*share.Link, error) {
	links, err := d.store.Share.All()
	if err != nil {
		return nil, err
	}

	adminPath := filepath.Clean(d.user.FullPath(path))
	owners := make(map[uint]*users.User)
	filtered := make([]*share.Link, 0, len(links))
	for _, link := range links {
		owner, ok := owners[link.UserID]
		if !ok {
			owner, err = d.store.Users.Get(d.server.Root, d.server.FollowExternalSymlinks, link.UserID)
			if err != nil && !errors.Is(err, fberrors.ErrNotExist) {
				return nil, err
			}
			owners[link.UserID] = owner // owner is nil on ErrNotExist
		}
		if owner != nil && filepath.Clean(owner.FullPath(link.Path)) == adminPath {
			filtered = append(filtered, link)
		}
	}

	return filtered, nil
}

// shareDeleteHandler serves DELETE /api/share/{hash}, for the link's owner or
// an admin.
var shareDeleteHandler = withPermShare(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
	hash := strings.TrimSuffix(r.URL.Path, "/")
	hash = strings.TrimPrefix(hash, "/")

	if hash == "" {
		return http.StatusBadRequest, nil
	}

	link, err := d.store.Share.GetByHash(hash)
	if err != nil {
		return errToStatus(err), err
	}

	if link.UserID != d.user.ID && !d.user.Perm.Admin {
		return http.StatusForbidden, nil
	}

	err = d.store.Share.Delete(hash)
	return errToStatus(err), err
})

// sharePostHandler serves POST /api/share/{path}, creating a public link with
// an optional password and expiry.
var sharePostHandler = withPermShare(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// Only allow sharing paths that currently exist. Otherwise a share could be
	// created for a non-existent path and would silently start exposing
	// whatever file later appears there.
	//
	// d.user.Fs is scoped, so Stat also refuses to follow a symlink whose target
	// escapes the user's scope: that returns a permission error here and so
	// blocks creating a share that points out of scope.
	if !d.Check(r.URL.Path) {
		return http.StatusForbidden, nil
	}
	if _, err := d.user.Fs.Stat(r.URL.Path); err != nil {
		return errToStatus(err), err
	}

	var s *share.Link
	var body share.CreateBody
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to decode body: %w", err)
		}
		defer r.Body.Close()
	}

	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	str := base64.URLEncoding.EncodeToString(bytes)

	expire, err := shareExpiry(body.Expires, body.Unit)
	if err != nil {
		return http.StatusBadRequest, err
	}

	hash, status, err := getSharePasswordHash(body)
	if err != nil {
		return status, err
	}

	var token string
	if len(hash) > 0 {
		tokenBuffer := make([]byte, 96)
		if _, err := rand.Read(tokenBuffer); err != nil {
			return http.StatusInternalServerError, err
		}
		token = base64.URLEncoding.EncodeToString(tokenBuffer)
	}

	s = &share.Link{
		Path:         r.URL.Path,
		Hash:         str,
		Expire:       expire,
		UserID:       d.user.ID,
		PasswordHash: string(hash),
		Token:        token,
	}

	if err := d.store.Share.Save(s); err != nil {
		return http.StatusInternalServerError, err
	}

	return renderJSON(w, r, toShareResponse(s))
})

// maxShareLifetime caps how far in the future a share may expire, which also
// keeps the duration arithmetic below from overflowing.
const maxShareLifetime = 100 * 365 * 24 * time.Hour

// shareExpiry turns the "expires" amount and unit of a create request into a
// Unix timestamp. An empty amount means the share never expires (0).
func shareExpiry(amount, unit string) (int64, error) {
	if amount == "" {
		return 0, nil
	}

	num, err := strconv.ParseInt(amount, 10, 64)
	if err != nil || num <= 0 {
		return 0, fmt.Errorf("invalid expiry %q: %w", amount, fberrors.ErrInvalidRequestParams)
	}

	var step time.Duration
	switch unit {
	case "seconds":
		step = time.Second
	case "minutes":
		step = time.Minute
	case "days":
		step = 24 * time.Hour
	default:
		step = time.Hour
	}

	if num > int64(maxShareLifetime/step) {
		return 0, fmt.Errorf("expiry %s %s is too far away: %w", amount, unit, fberrors.ErrInvalidRequestParams)
	}

	return time.Now().Add(time.Duration(num) * step).Unix(), nil
}

// getSharePasswordHash bcrypt-hashes the share password, if one was given.
func getSharePasswordHash(body share.CreateBody) (data []byte, statuscode int, err error) {
	if body.Password == "" {
		return nil, 0, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, 0, nil
}
