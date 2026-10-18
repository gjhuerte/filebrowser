package bolt

import (
	"errors"

	"github.com/asdine/storm/v3"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
)

// get reads the value stored under name in the "config" bucket into to.
func get(db *storm.DB, name string, to any) error {
	err := db.Get("config", name, to)
	if errors.Is(err, storm.ErrNotFound) {
		return fberrors.ErrNotExist
	}

	return err
}

// save stores from under name in the "config" bucket.
func save(db *storm.DB, name string, from any) error {
	return db.Set("config", name, from)
}
