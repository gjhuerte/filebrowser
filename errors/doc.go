// Package fberrors defines the error values shared across File Browser.
//
// HTTP handlers map them to status codes (for example ErrNotExist to 404), so
// callers should compare with errors.Is rather than by message.
package fberrors
