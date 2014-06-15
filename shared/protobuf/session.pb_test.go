package protobuf

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSession_Reset(t *testing.T) {
	expiresAt := time.Now().Unix()
	session := Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	session.Reset()
	if session.Values != nil || session.ExpiresAt != nil {
		t.Error("the session should be zero value")
	}
}

func TestSession_String(t *testing.T) {
	expiresAt := time.Now().Unix()
	session := Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	expected := fmt.Sprintf(`Values:"test" ExpiresAt:%d`, expiresAt)
	actual := strings.TrimSpace(session.String())
	if actual != expected {
		t.Errorf("session.String() should return %s (actual: %s)", expected, actual)
	}
}

func TestSession_GetValues(t *testing.T) {
	// When session == nil.
	var session *Session
	actual := session.GetValues()
	if actual != nil {
		t.Errorf("session.GetValues() should return nil (actual: %+v)", actual)
	}

	// When session != nil.
	expiresAt := time.Now().Unix()
	session = &Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	expected := []byte("test")
	actual = session.GetValues()
	if len(actual) != len(expected) || actual[0] != expected[0] {
		t.Errorf("session.GetValues() should return %+v (actual: %+v)", expected, actual)
	}

	var s *Session
	s.GetValues()
}

func TestSession_GetExpiresAt(t *testing.T) {
	// When Session.ExpiresAt == nil.
	session := Session{
		Values:    []byte("test"),
		ExpiresAt: nil,
	}
	expected := int64(0)
	actual := session.GetExpiresAt()
	if actual != expected {
		t.Errorf("session.GetExpiresAt() should return %d (actual: %d)", expected, actual)
	}

	// When Session.ExpiresAt != nil.
	expiresAt := time.Now().Unix()
	session = Session{
		Values:    []byte("test"),
		ExpiresAt: &expiresAt,
	}
	expected = expiresAt
	actual = session.GetExpiresAt()
	if actual != expected {
		t.Errorf("session.GetExpiresAt() should return %d (actual: %d)", expected, actual)
	}
	session.ProtoMessage()
}
