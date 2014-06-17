package store

import (
	"github.com/gorilla/sessions"
	"github.com/yosssi/boltstore/shared"
)

// Config represents a config for a session store.
type Config struct {
	SessionOptions sessions.Options
	DBOptions      Options
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
