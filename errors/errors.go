package fberrors

import (
	"errors"
	"fmt"
)

// Errors returned across File Browser. The HTTP status each one maps to is
// noted where handlers translate it (see errToStatus in package fbhttp).
var (
	// ErrEmptyKey means the settings have no signing key for auth tokens.
	ErrEmptyKey = errors.New("empty key")
	// ErrExist means the resource already exists (409 Conflict).
	ErrExist = errors.New("the resource already exists")
	// ErrNotExist means the resource does not exist (404 Not Found).
	ErrNotExist = errors.New("the resource does not exist")
	// ErrEmptyPassword means a user was saved without a password.
	ErrEmptyPassword = errors.New("password is empty")
	// ErrEasyPassword means the password is on the list of common passwords.
	ErrEasyPassword = errors.New("password is too easy")
	// ErrEmptyUsername means a user was saved without a username.
	ErrEmptyUsername = errors.New("username is empty")
	// ErrEmptyRequest means a request body that is required was missing.
	ErrEmptyRequest = errors.New("empty request")
	// ErrInvalidDataType means a value had a type the function cannot handle,
	// e.g. a user lookup key that is neither a uint ID nor a username.
	ErrInvalidDataType = errors.New("invalid data type")
	// ErrIsDirectory means a file operation was attempted on a folder
	// (409 Conflict).
	ErrIsDirectory = errors.New("file is directory")
	// ErrInvalidOption means an option value is not supported, e.g. an
	// unknown checksum algorithm.
	ErrInvalidOption = errors.New("invalid option")
	// ErrInvalidAuthMethod means the configured auth method is unknown.
	ErrInvalidAuthMethod = errors.New("invalid auth method")
	// ErrPermissionDenied means the user's permissions or rules forbid the
	// operation (403 Forbidden).
	ErrPermissionDenied = errors.New("permission denied")
	// ErrInvalidRequestParams means the request parameters are malformed
	// (400 Bad Request).
	ErrInvalidRequestParams = errors.New("invalid request params")
	// ErrSourceIsParent means a copy or move would put a folder inside itself.
	ErrSourceIsParent = errors.New("source is parent")
	// ErrRootUserDeletion means the request would delete the last admin.
	ErrRootUserDeletion = errors.New("the sole admin can't be deleted")
	// ErrCurrentPasswordIncorrect means the confirmation password for a
	// sensitive change did not match.
	ErrCurrentPasswordIncorrect = errors.New("the current password is incorrect")
	// ErrShareRequiresDownload means a user was given the Share permission
	// without Download, which sharing depends on.
	ErrShareRequiresDownload = errors.New("permission to share requires permission to download")
)

// ErrShortPassword means the password is shorter than the configured minimum.
type ErrShortPassword struct {
	// MinimumLength is the shortest password accepted.
	MinimumLength uint
}

// Error implements the error interface.
func (e ErrShortPassword) Error() string {
	return fmt.Sprintf("password is too short, minimum length is %d", e.MinimumLength)
}
