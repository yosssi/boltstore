package store

import (
	"net/http"
	"testing"

	"github.com/boltdb/bolt"
)

func TestStore_Get(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	str, err := New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	if session.IsNew != true {
		t.Errorf("session.IsNew should be true (actual: %+v)", session.IsNew)
	}
}
