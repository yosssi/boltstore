package reaper

import (
	"testing"

	"github.com/yosssi/boltstore/shared"
)

func TestOptions_setDefault(t *testing.T) {
	options := Options{}
	options.setDefault()
	if string(options.BucketName) != shared.DefaultBucketName {
		t.Errorf("options.BucketName should be %+v (actual: %+v)", []byte(shared.DefaultBucketName), options.BucketName)
	}
	if options.BatchSize != shared.DefaultBatchSize {
		t.Errorf("options.BucketName should be %+d (actual: %+d)", shared.DefaultBatchSize, options.BatchSize)
	}
	if options.CheckInterval != shared.DefaultCheckInterval {
		t.Errorf("options.BucketName should be %+v (actual: %+v)", shared.DefaultCheckInterval, options.CheckInterval)
	}
}
