package store

import (
	"github.com/gorilla/sessions"
	"github.com/yosssi/boltstore/shared"
)

// Config represents a config for a session store.
type Config struct {
	// SessionOptions represents options for a session.
	SessionOptions sessions.Options
	// DBOptions represents options for a database.
	DBOptions Options
}

// setDefault sets default to the config.
func (c *Config) setDefault() {
	if c.SessionOptions.Path == "" {
		c.SessionOptions.Path = shared.DefaultPath
	}
	if c.SessionOptions.MaxAge == 0 {
		c.SessionOptions.MaxAge = shared.DefaultMaxAge
	}
	if c.DBOptions.BucketName == nil {
		c.DBOptions.BucketName = []byte(shared.DefaultBucketName)
	}
}
