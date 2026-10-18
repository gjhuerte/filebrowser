package share

import (
	"time"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
)

// StorageBackend is the interface to implement for a share storage.
type StorageBackend interface {
	All() ([]*Link, error)
	FindByUserID(id uint) ([]*Link, error)
	GetByHash(hash string) (*Link, error)
	Gets(path string, id uint) ([]*Link, error)
	Save(s *Link) error
	Delete(hash string) error
	DeleteWithPathPrefix(path string, userID uint) error
}

// Storage is the share links store. Expired links are deleted lazily, the
// first time a lookup comes across them.
type Storage struct {
	back StorageBackend
}

// NewStorage creates a share links storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

// expired reports whether the link has an expiry date that has passed.
func (l *Link) expired(now int64) bool {
	return l.Expire != 0 && l.Expire <= now
}

// dropExpired deletes the expired links from the backend and returns the rest.
func (s *Storage) dropExpired(links []*Link) ([]*Link, error) {
	now := time.Now().Unix()
	valid := links[:0]
	for _, link := range links {
		if !link.expired(now) {
			valid = append(valid, link)
			continue
		}
		if err := s.Delete(link.Hash); err != nil {
			return nil, err
		}
	}
	return valid, nil
}

// All returns every link that has not expired.
func (s *Storage) All() ([]*Link, error) {
	links, err := s.back.All()
	if err != nil {
		return nil, err
	}

	return s.dropExpired(links)
}

// FindByUserID returns the unexpired links created by the given user.
func (s *Storage) FindByUserID(id uint) ([]*Link, error) {
	links, err := s.back.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	return s.dropExpired(links)
}

// GetByHash returns the link with the given hash, or fberrors.ErrNotExist if
// there is none or it has expired.
func (s *Storage) GetByHash(hash string) (*Link, error) {
	link, err := s.back.GetByHash(hash)
	if err != nil {
		return nil, err
	}

	if link.expired(time.Now().Unix()) {
		if err := s.Delete(link.Hash); err != nil {
			return nil, err
		}
		return nil, fberrors.ErrNotExist
	}

	return link, nil
}

// Gets returns the user's unexpired links for path.
func (s *Storage) Gets(path string, id uint) ([]*Link, error) {
	links, err := s.back.Gets(path, id)
	if err != nil {
		return nil, err
	}

	return s.dropExpired(links)
}

// Save stores a link, replacing any link with the same hash.
func (s *Storage) Save(l *Link) error {
	return s.back.Save(l)
}

// Delete removes the link with the given hash.
func (s *Storage) Delete(hash string) error {
	return s.back.Delete(hash)
}

// DeleteWithPathPrefix removes the user's links for path and everything below
// it, e.g. after the folder was deleted.
func (s *Storage) DeleteWithPathPrefix(path string, userID uint) error {
	return s.back.DeleteWithPathPrefix(path, userID)
}
