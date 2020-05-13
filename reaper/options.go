package reaper

import (
	"time"

	"github.com/yosssi/boltstore/shared"
)

// Options represents options for the reaper.
type Options struct {
	// BucketName represents the name of the bucket which contains sessions.
	BucketName []byte
	// BatchSize represents the maximum number of sessions which the reaper
	// process at one time.
	BatchSize int
	// CheckInterval represents the interval between the reaper's invocation.
	CheckInterval time.Duration
	// PreDeleteFn register a function to be called before deleting a session.
	// If a non-nil error is returned, the session won't be deleted
	PreDeleteFn func(values map[interface{}]interface{}) error
}

// setDefault sets default to the reaper options.
func (o *Options) setDefault() {
	if o.BucketName == nil {
		o.BucketName = []byte(shared.DefaultBucketName)
	}
	if o.BatchSize == 0 {
		o.BatchSize = shared.DefaultBatchSize
	}
	if o.CheckInterval == 0 {
		o.CheckInterval = shared.DefaultCheckInterval
	}
}
