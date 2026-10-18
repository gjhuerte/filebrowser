package bolt

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	bolt "go.etcd.io/bbolt"

	fberrors "github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/users"
)

// usersBackend stores users, indexed by ID and unique username.
type usersBackend struct {
	db *storm.DB
}

// GetBy implements users.StorageBackend: i is a uint ID or a string username.
func (st usersBackend) GetBy(i any) (user *users.User, err error) {
	user = &users.User{}

	var arg string
	switch i.(type) {
	case uint:
		arg = "ID"
	case string:
		arg = "Username"
	default:
		return nil, fberrors.ErrInvalidDataType
	}

	err = st.db.One(arg, i, user)

	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil, fberrors.ErrNotExist
		}
		return nil, err
	}

	return
}

// GetByScope implements users.StorageBackend.
func (st usersBackend) GetByScope(scope string) (*users.User, error) {
	user := &users.User{}
	// Match case-insensitively: on a case-insensitive filesystem two scopes
	// that differ only in case (e.g. /users/Alice and /users/alice) resolve to
	// the same home directory, so they must be treated as a collision.
	pattern := "(?i)^" + regexp.QuoteMeta(scope) + "$"
	err := st.db.Select(q.Re("Scope", pattern)).First(user)
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil, fberrors.ErrNotExist
		}
		return nil, err
	}

	return user, nil
}

// Gets implements users.StorageBackend.
func (st usersBackend) Gets() ([]*users.User, error) {
	var allUsers []*users.User
	err := st.db.All(&allUsers)
	if errors.Is(err, storm.ErrNotFound) {
		return nil, fberrors.ErrNotExist
	}

	if err != nil {
		return allUsers, err
	}

	return allUsers, err
}

// Update implements users.StorageBackend, writing only the named fields, or
// the whole user when none are given.
func (st usersBackend) Update(user *users.User, fields ...string) error {
	if len(fields) == 0 {
		return st.Save(user)
	}

	for _, field := range fields {
		userField := reflect.ValueOf(user).Elem().FieldByName(field)
		if !userField.IsValid() {
			return fmt.Errorf("invalid field: %s", field)
		}
		val := userField.Interface()
		if err := st.db.UpdateField(user, field, val); err != nil {
			return err
		}
	}

	return nil
}

// Save implements users.StorageBackend.
func (st usersBackend) Save(user *users.User) error {
	err := st.db.Save(user)
	if errors.Is(err, storm.ErrAlreadyExists) {
		return fberrors.ErrExist
	}
	return err
}

// DeleteByID implements users.StorageBackend.
func (st usersBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&users.User{ID: id})
}

// DeleteByUsername implements users.StorageBackend.
func (st usersBackend) DeleteByUsername(username string) error {
	user, err := st.GetBy(username)
	if err != nil {
		return err
	}

	return st.db.DeleteStruct(user)
}

// CountAdmins implements users.StorageBackend.
func (st usersBackend) CountAdmins() (int, error) {
	count := 0

	err := st.db.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(reflect.TypeOf(users.User{}).Name()))
		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()
		for _, v := c.First(); v != nil; _, v = c.Next() {
			var u users.User
			if err := st.db.Codec().Unmarshal(v, &u); err != nil {
				return err
			}
			if u.Perm.Admin {
				count++
			}
		}

		return nil
	})

	return count, err
}
