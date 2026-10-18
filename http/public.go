package fbhttp

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/share"
	"golang.org/x/crypto/bcrypt"
)

// withHashFile resolves a public share from the link hash in the path, checks
// its password (with rate limiting) and rebases the owner's filesystem onto the
// shared item, so the rest of the request can only see what was shared.
var withHashFile = func(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, ifPath := ifPathWithName(r)
		link, err := d.store.Share.GetByHash(id)
		if err != nil {
			return errToStatus(err), err
		}

		// Password-protected shares are public, so limit guessing per client
		// and share. A request without a password (the first visit, before
		// the prompt) is not counted as a failure.
		limitKey := clientKey(r) + "|" + link.Hash
		if link.PasswordHash != "" {
			if ok, wait := shareLimiter.Allow(limitKey); !ok {
				return tooManyAttempts(w, wait), nil
			}
		}

		status, err := authenticateShareRequest(r, link)
		if status == http.StatusUnauthorized && r.Header.Get("X-SHARE-PASSWORD") != "" {
			shareLimiter.Fail(limitKey)
		}
		if status != 0 || err != nil {
			return status, err
		}
		if link.PasswordHash != "" {
			shareLimiter.Reset(limitKey)
		}

		user, err := d.store.Users.Get(d.server.Root, d.server.FollowExternalSymlinks, link.UserID)
		if err != nil {
			return errToStatus(err), err
		}

		if !user.Perm.Share || !user.Perm.Download {
			return http.StatusForbidden, nil
		}

		d.user = user

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       link.Path,
			Modify:     d.user.Perm.Modify,
			Expand:     false,
			ReadHeader: d.server.TypeDetectionByHeader,
			CalcImgRes: d.server.TypeDetectionByHeader,
			Checker:    d,
			Token:      link.Token,
		})
		if err != nil {
			return errToStatus(err), err
		}

		// share base path. Canonicalized because it roots both the rebased
		// filesystem and checkerPrefix below, and a stored path that is not
		// "/"-separated would make the two disagree on Windows.
		basePath := slashClean(link.Path)

		// file relative path
		filePath := ""

		if file.IsDir {
			filePath = ifPath
		}

		// set fs root to the shared file/folder. Unless external symlinks are
		// explicitly allowed, this is a ScopedFs (not a bare BasePathFs) so the
		// share is also symlink-confined: a link inside the shared subtree that
		// points elsewhere in the owner's scope — outside the share — must not be
		// followed.
		d.user.Fs = files.NewFs(d.user.Fs, basePath, d.server.FollowExternalSymlinks)

		// the filesystem is now rebased onto basePath, so paths handed to the
		// rule checker are relative to it. Resolve them back to the user's
		// original scope so deny rules below the share root keep applying.
		d.checkerPrefix = basePath

		file, err = files.NewFileInfo(&files.FileOptions{
			Fs:      d.user.Fs,
			Path:    filePath,
			Modify:  d.user.Perm.Modify,
			Expand:  true,
			Checker: d,
			Token:   link.Token,
		})
		if err != nil {
			return errToStatus(err), err
		}

		if file.IsDir {
			// extract name from the last directory in the path
			name := filepath.Base(strings.TrimRight(link.Path, string(filepath.Separator)))
			file.Name = name
		}

		d.raw = file
		return fn(w, r, d)
	}
}

// ref to https://github.com/filebrowser/filebrowser/pull/727
// `/api/public/dl/MEEuZK-v/file-name.txt` for old browsers to save file with correct name
func ifPathWithName(r *http.Request) (id, filePath string) {
	pathElements := strings.Split(r.URL.Path, "/")
	// prevent maliciously constructed parameters like `/api/public/dl/XZzCDnK2_not_exists_hash_name`
	// len(pathElements) will be 1, and golang will panic `runtime error: index out of range`

	switch len(pathElements) {
	case 1:
		return r.URL.Path, "/"
	default:
		// Public share routes do not pass through withUser, so canonicalize the
		// share-relative path here instead.
		return pathElements[0], slashClean(path.Join(pathElements[1:]...))
	}
}

// publicShareHandler serves GET /api/public/share/{hash}/{path}: details or a
// listing of a shared item.
var publicShareHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)

	if file.IsDir {
		file.Sorting = files.Sorting{By: "name", Asc: false}
		file.ApplySort()
		return renderJSON(w, r, file)
	}

	return renderJSON(w, r, file)
})

// publicDlHandler serves GET /api/public/dl/{hash}/{path}: the shared file, or
// a shared folder as an archive.
var publicDlHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)
	if !file.IsDir {
		return rawFileHandler(w, r, file)
	}

	return rawDirHandler(w, r, d, file)
})

// authenticateShareRequest checks a password-protected share: either the
// share's download token (?token=) or the X-SHARE-PASSWORD header. It returns
// 401 when neither matches.
func authenticateShareRequest(r *http.Request, l *share.Link) (int, error) {
	if l.PasswordHash == "" {
		return 0, nil
	}

	if subtle.ConstantTimeCompare([]byte(r.URL.Query().Get("token")), []byte(l.Token)) == 1 {
		return 0, nil
	}

	password := r.Header.Get("X-SHARE-PASSWORD")
	password, err := url.QueryUnescape(password)
	if err != nil {
		return 0, err
	}
	if password == "" {
		return http.StatusUnauthorized, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(l.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return http.StatusUnauthorized, nil
		}
		return 0, err
	}

	return 0, nil
}

// healthHandler serves GET /health for container health checks.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"OK"}`))
}
