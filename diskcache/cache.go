package diskcache

import (
	"context"
)

// Interface is a byte cache keyed by string. Load reports whether the key
// was found; a missing key is not an error.
type Interface interface {
	Store(ctx context.Context, key string, value []byte) error
	Load(ctx context.Context, key string) (value []byte, exist bool, err error)
	Delete(ctx context.Context, key string) error
}
