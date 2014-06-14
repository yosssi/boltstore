package reaper

import (
	"time"

	"github.com/yosssi/boltstore/shared"
)

// Options represents options for the reaper.
type Options struct {
	BucketName    []byte
	BatchSize     int
	CheckInterval time.Duration
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
