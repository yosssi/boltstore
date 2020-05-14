package store

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"github.com/gogo/protobuf/proto"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/yosssi/boltstore/shared"
)

// Store represents a session store.
type Store struct {
	codecs []securecookie.Codec
	config Config
	db     *bolt.DB
}

// Get returns a session for the given name after adding it to the registry.
//
// See gorilla/sessions FilesystemStore.Get().
func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See gorilla/sessions FilesystemStore.New().
func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	var err error
	session := sessions.NewSession(s, name)
	session.Options = &s.config.SessionOptions
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.codecs...)
		if err == nil {
			ok, err := s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		s.delete(session)
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric ID.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// load loads a session data from the database.
// True is returned if there is a session data in the database.
func (s *Store) load(session *sessions.Session) (bool, error) {
	// exists represents whether a session data exists or not.
	var exists bool
	var expired bool
	id := []byte(session.ID)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.config.DBOptions.BucketName)
		// Get the session data.
		data := bucket.Get(id)
		if data == nil {
			return nil
		}
		sessionData, err := shared.Session(data)
		if err != nil {
			return err
		}
		// Check the expiration of the session data.
		if shared.Expired(sessionData) {
			expired = true
			return err
		}
		exists = true
		dec := gob.NewDecoder(bytes.NewBuffer(sessionData.Values))
		return dec.Decode(&session.Values)
	})
	if expired {
		if updateErr := s.db.Update(func(txu *bolt.Tx) error {
			return txu.Bucket(s.config.DBOptions.BucketName).Delete(id)
		}); updateErr != nil {
			err = updateErr
		}
	}
	return exists, err
}

// delete removes the key-value from the database.
func (s *Store) delete(session *sessions.Session) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.config.DBOptions.BucketName).Delete([]byte(session.ID))
	})
	if err != nil {
		return err
	}
	return nil
}

// save stores the session data in the database.
func (s *Store) save(session *sessions.Session) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(session.Values)
	if err != nil {
		return err
	}
	data, err := proto.Marshal(shared.NewSession(buf.Bytes(), session.Options.MaxAge))
	if err != nil {
		return err
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.config.DBOptions.BucketName).Put([]byte(session.ID), data)
	})
	return err
}

// New creates and returns a session store.
func New(db *bolt.DB, config Config, keyPairs ...[]byte) (*Store, error) {
	config.setDefault()
	store := &Store{
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		config: config,
		db:     db,
	}
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(config.DBOptions.BucketName)
		return err
	})
	if err != nil {
		return nil, err
	}
	return store, nil
}
