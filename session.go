package boltstore

import "time"

// NewSession creates and returns a session data.
func NewSession(values []byte, maxAge int) *Session {
	var expiresAt int64
	if maxAge > 0 {
		expiresAt = time.Now().Unix() + int64(maxAge)
	}
	return &Session{Values: values, ExpiresAt: &expiresAt}
}
