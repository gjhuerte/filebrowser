package share

import (
	"testing"
	"time"
)

// memoryBackend is a minimal in-memory StorageBackend for tests.
type memoryBackend struct {
	links   []*Link
	deleted []string
}

func (m *memoryBackend) All() ([]*Link, error) {
	return append([]*Link(nil), m.links...), nil
}
func (m *memoryBackend) FindByUserID(uint) ([]*Link, error)      { return m.All() }
func (m *memoryBackend) GetByHash(string) (*Link, error)         { return m.links[0], nil }
func (m *memoryBackend) Gets(string, uint) ([]*Link, error)      { return m.All() }
func (m *memoryBackend) Save(*Link) error                        { return nil }
func (m *memoryBackend) DeleteWithPathPrefix(string, uint) error { return nil }
func (m *memoryBackend) Delete(hash string) error {
	m.deleted = append(m.deleted, hash)
	return nil
}

// Consecutive expired links used to be skipped (or to panic) because the slice
// was modified while it was being ranged over.
func TestAllDropsEveryExpiredLink(t *testing.T) {
	past := time.Now().Add(-time.Hour).Unix()
	future := time.Now().Add(time.Hour).Unix()
	back := &memoryBackend{links: []*Link{
		{Hash: "keep-1"},
		{Hash: "old-1", Expire: past},
		{Hash: "old-2", Expire: past},
		{Hash: "keep-2", Expire: future},
		{Hash: "old-3", Expire: past},
	}}

	links, err := NewStorage(back).All()
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, l := range links {
		got = append(got, l.Hash)
	}
	if len(got) != 2 || got[0] != "keep-1" || got[1] != "keep-2" {
		t.Errorf("unexpected links after filtering: %v", got)
	}
	if len(back.deleted) != 3 {
		t.Errorf("expected 3 expired links deleted, got %v", back.deleted)
	}
}
