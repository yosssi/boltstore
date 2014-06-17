package shared

import (
	"testing"
	"time"
	"github.com/yosssi/boltstore/shared/protobuf"

	"code.google.com/p/gogoprotobuf/proto"
)

func TestSession(t *testing.T) {
	expiresAt := time.Now().Unix()
	sessionOrig := &protobuf.Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	data, err := proto.Marshal(sessionOrig)
	if err != nil {
		t.Error(err.Error())
	}
	session, err := Session(data)
	if err != nil {
		t.Error(err.Error())
	}
	if string(session.Values) != "test" || *session.ExpiresAt != expiresAt {
		t.Errorf("Session() should return %+v (actual: %+v)", sessionOrig, session)
	}
}

func TestExpired(t *testing.T) {
	// When Expired() should return true.
	expiresAt := time.Now().Unix()
	session := protobuf.Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	if Expired(session) != true {
		t.Error("Expired() should return true (actual: false)")
	}

	// When Expired() should return false.
	expiresAt = time.Now().Add(time.Hour).Unix()
	session = protobuf.Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	if Expired(session) != false {
		t.Error("Expired() should return false (actual: true)")
	}
}

func TestNewSession(t *testing.T) {
	values := []byte("test")
	maxAge := 10
	preExpiresAt := time.Now().Unix() + int64(maxAge)
	session := NewSession(values, maxAge)
	postExpiresAt := time.Now().Unix() + int64(maxAge)
	if string(session.Values) != string(values) || *session.ExpiresAt < preExpiresAt || postExpiresAt < *session.ExpiresAt {
		t.Errorf("NewSession() returned an invalid value (actual: %+v)", session)
	}
}
