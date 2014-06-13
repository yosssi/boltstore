package boltstore

import "time"

const (
	defaultPath          = "/"
	defaultDBPath        = "./sessions.db"
	defaultBucketName    = "sessions"
	defaultBatchSize     = 1000
	defaultCheckInterval = 10 * time.Second
)
