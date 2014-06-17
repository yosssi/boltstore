package store

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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

func TestStore_New(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	str, err := New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	encoded, err := securecookie.EncodeMulti("test", "1", str.codecs...)

	req.AddCookie(sessions.NewCookie("test", encoded, &sessions.Options{
		MaxAge: 1024,
	}))

	session, err := str.New(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	if session.IsNew != true {
		t.Errorf("session.IsNew should be true (actual: %+v)", session.IsNew)
	}
}

func TestStore_Save(t *testing.T) {
	// When session.Options.MaxAge < 0
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	str, err := New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	session.Options.MaxAge = -1

	w := httptest.NewRecorder()

	str.Save(req, w, session)

	// When session.Options.MaxAge >= 0
	session, err = str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	session.Options.MaxAge = 1

	w = httptest.NewRecorder()

	str.Save(req, w, session)

	// When session.Options.MaxAge >= 0 and
	// s.save returns an error
	session.Values = make(map[interface{}]interface{})
	session.Values[make(chan int)] = make(chan int)
	str.Save(req, w, session)

	// When session.Options.MaxAge >= 0 and
	// securecookie.EncodeMulti  returns an error
	session.Values = make(map[interface{}]interface{})
	str.codecs = nil
	str.Save(req, w, session)
}

func TestStore_load(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	str, err := New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	w := httptest.NewRecorder()

	str.Save(req, w, session)

	exists, err := str.load(session)
	if err != nil {
		t.Error(err.Error())
	}

	if exists != true {
		t.Error("Store.load should return true (actual: false)")
	}

	// When the target session data is nil
	session.ID = "x"
	exists, err = str.load(session)
	if err != nil {
		t.Error(err.Error())
	}

	if exists != false {
		t.Error("Store.load should return false (actual: true)")
	}

	// When shared.Session returns an error
	err = db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(str.config.DBOptions.BucketName).Put([]byte("x"), []byte("test"))
	})
	if err != nil {
		t.Error(err.Error())
	}
	_, err = str.load(session)
	if err == nil || err.Error() != "proto: field/encoding mismatch: wrong type for field" {
		t.Error(`str.load should return an error "%s" (actual: %s)`, "proto: field/encoding mismatch: wrong type for field", err)
	}

	// When the target session data is expired
	session, err = str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}
	session.Options.MaxAge = 0
	str.Save(req, w, session)
	time.Sleep(time.Second)
	_, err = str.load(session)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSession_delete(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	str, err := New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	session, err := str.Get(req, "test")
	if err != nil {
		t.Error(err.Error())
	}

	db.Close()

	err = str.delete(session)

	if err.Error() != "database not open" {
		t.Error(`str.delete  should return an error "%s" (actual: %s)`, "database not open", err)
	}
}

func TestNew(t *testing.T) {
	// When db.Update returns an error
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	db.Close()

	_, err = New(
		db,
		Config{},
		[]byte("secret-key"),
	)
	if err.Error() != "database not open" {
		t.Error(`str.delete  should return an error "%s" (actual: %s)`, "database not open", err)
	}
}
