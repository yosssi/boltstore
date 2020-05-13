package store

import (
	"testing"

	"github.com/yosssi/boltstore/shared"
)

func TestConfig_setDefault(t *testing.T) {
	config := Config{}
	config.setDefault()
	if config.SessionOptions.Path != shared.DefaultPath {
		t.Errorf("config.SessionOptions.Path should be %s (actual: %s)", shared.DefaultPath, config.SessionOptions.Path)
	}
	if config.SessionOptions.MaxAge != shared.DefaultMaxAge {
		t.Errorf("config.SessionOptions.MaxAge should be %d (actual: %d)", shared.DefaultMaxAge, config.SessionOptions.MaxAge)
	}
	if string(config.DBOptions.BucketName) != shared.DefaultBucketName {
		t.Errorf("config.SessionOptions.BucketName should be %+v (actual: %+v)", shared.DefaultBucketName, config.DBOptions.BucketName)
	}
}
