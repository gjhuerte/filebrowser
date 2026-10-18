package diskcache

import (
	"context"
)

// NoOp is an Interface that stores nothing. It is used when previews are not
// cached.
type NoOp struct {
}

// NewNoOp returns a cache that discards everything.
func NewNoOp() *NoOp {
	return &NoOp{}
}

// Store does nothing.
func (n *NoOp) Store(_ context.Context, _ string, _ []byte) error {
	return nil
}

// Load always reports that the key does not exist.
func (n *NoOp) Load(_ context.Context, _ string) (value []byte, exist bool, err error) {
	return nil, false, nil
}

// Delete does nothing.
func (n *NoOp) Delete(_ context.Context, _ string) error {
	return nil
}
