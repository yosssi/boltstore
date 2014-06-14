package store

import (
	"time"

	"github.com/yosssi/boltstore/shared/protobuf"
)

// NewSession creates and returns a session data.
func NewSession(values []byte, maxAge int) *protobuf.Session {
	var expiresAt int64
	if maxAge > 0 {
		expiresAt = time.Now().Unix() + int64(maxAge)
	}
	return &protobuf.Session{Values: values, ExpiresAt: &expiresAt}
}
