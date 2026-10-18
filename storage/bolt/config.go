package bolt

import (
	"github.com/asdine/storm/v3"

	"github.com/filebrowser/filebrowser/v2/settings"
)

// settingsBackend stores the settings and server options as two keys in the
// "config" bucket.
type settingsBackend struct {
	db *storm.DB
}

// Get implements settings.StorageBackend.
func (s settingsBackend) Get() (*settings.Settings, error) {
	set := &settings.Settings{}
	return set, get(s.db, "settings", set)
}

// Save implements settings.StorageBackend.
func (s settingsBackend) Save(set *settings.Settings) error {
	return save(s.db, "settings", set)
}

// GetServer implements settings.StorageBackend.
func (s settingsBackend) GetServer() (*settings.Server, error) {
	server := &settings.Server{}
	return server, get(s.db, "server", server)
}

// SaveServer implements settings.StorageBackend.
func (s settingsBackend) SaveServer(server *settings.Server) error {
	return save(s.db, "server", server)
}
