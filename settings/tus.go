package settings

// Defaults for resumable uploads.
const (
	// DefaultTusChunkSize is the size of each upload chunk, 10 MiB.
	DefaultTusChunkSize = 10 * 1024 * 1024
	// DefaultTusRetryCount is how often the frontend retries a failed chunk.
	DefaultTusRetryCount = 5
)

// Tus contains the tus.io settings of the app.
type Tus struct {
	ChunkSize  uint64 `json:"chunkSize"`
	RetryCount uint16 `json:"retryCount"`
}
