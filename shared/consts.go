package shared

import "time"

// Defaults for store.Options
const (
	DefaultPath       = "/"
	DefaultBucketName = "sessions"
)

// Defaults for reaper.Options
const (
	DefaultBatchSize     = 10
	DefaultCheckInterval = time.Second
)
