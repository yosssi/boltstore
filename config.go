package boltstore

import "github.com/gorilla/sessions"

// Config represents a config for a session store.
type Config struct {
	SessionOptions sessions.Options
	DBOptions      Options
}

// setDefault sets default to the config.
func (c *Config) setDefault() {
	if c.SessionOptions.Path == "" {
		c.SessionOptions.Path = defaultPath
	}
	if c.DBOptions.Path == "" {
		c.DBOptions.Path = defaultDBPath
	}
	if c.DBOptions.BucketName == nil {
		c.DBOptions.BucketName = []byte(defaultBucketName)
	}
}
