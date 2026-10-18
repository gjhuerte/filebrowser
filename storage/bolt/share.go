package bolt

import (
	"errors"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/share"
)

// shareBackend stores share links, indexed by hash and path.
type shareBackend struct {
	db *storm.DB
}

// All implements share.StorageBackend.
func (s shareBackend) All() ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.All(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return v, fberrors.ErrNotExist
	}

	return v, err
}

// FindByUserID implements share.StorageBackend.
func (s shareBackend) FindByUserID(id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Select(q.Eq("UserID", id)).Find(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return v, fberrors.ErrNotExist
	}

	return v, err
}

// GetByHash implements share.StorageBackend.
func (s shareBackend) GetByHash(hash string) (*share.Link, error) {
	var v share.Link
	err := s.db.One("Hash", hash, &v)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, fberrors.ErrNotExist
	}

	return &v, err
}

// Gets implements share.StorageBackend.
func (s shareBackend) Gets(path string, id uint) ([]*share.Link, error) {
	var v []*share.Link
	err := s.db.Select(q.Eq("Path", path), q.Eq("UserID", id)).Find(&v)
	if errors.Is(err, storm.ErrNotFound) {
		return v, fberrors.ErrNotExist
	}

	return v, err
}

// Save implements share.StorageBackend.
func (s shareBackend) Save(l *share.Link) error {
	return s.db.Save(l)
}

// Delete implements share.StorageBackend. Deleting a missing link is not an
// error.
func (s shareBackend) Delete(hash string) error {
	err := s.db.DeleteStruct(&share.Link{Hash: hash})
	if errors.Is(err, storm.ErrNotFound) {
		return nil
	}
	return err
}

// DeleteWithPathPrefix implements share.StorageBackend.
func (s shareBackend) DeleteWithPathPrefix(pathPrefix string, userID uint) error {
	// Share paths are stored without a trailing slash
	prefix := strings.TrimRight(pathPrefix, "/")

	var links []share.Link
	if err := s.db.Prefix("Path", prefix, &links); err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil
		}
		return err
	}

	var err error
	for _, link := range links {
		if link.UserID != userID {
			continue
		}

		if link.Path != prefix && !strings.HasPrefix(link.Path, prefix+"/") {
			continue
		}

		err = errors.Join(err, s.db.DeleteStruct(&share.Link{Hash: link.Hash}))
	}
	return err
}
