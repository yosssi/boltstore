package reaper

import (
	"fmt"
	"testing"
	"time"
	"code.google.com/p/gogoprotobuf/proto"

	"github.com/boltdb/bolt"
	"github.com/yosssi/boltstore/shared"
)

func TestRun(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
	defer Quit(Run(db, Options{}))
}

func TestQuit(t *testing.T) {
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
	quitC, doneC := Run(db, Options{})
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)
}

func Test_reap(t *testing.T) {
	// When the target bucket does not exist
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
	options := Options{}
	options.setDefault()
	quitC, doneC := make(chan struct{}), make(chan struct{})
	go reap(db, options, quitC, doneC)
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)

	// When no keys exist
	bucketName := []byte(fmt.Sprintf("reapTest-%d", time.Now().Unix()))
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		t.Error(err.Error())
	}

	options.BucketName = bucketName

	go reap(db, options, quitC, doneC)
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)

	// When shared.Session returns an error
	err = db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte("test"), []byte("value"))
	})
	if err != nil {
		t.Error(err.Error())
	}
	go reap(db, options, quitC, doneC)
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)

	// When the target session is expired
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketName).Delete([]byte("test"))
		if err != nil {
			return err
		}
		session := shared.NewSession([]byte{}, -1)
		data, err := proto.Marshal(session)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketName).Put([]byte("test"), data)
	})
	if err != nil {
		t.Error(err.Error())
	}
	go reap(db, options, quitC, doneC)
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)

	// When options.BatchSize == i
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketName).Delete([]byte("test"))
		if err != nil {
			return err
		}
		session := shared.NewSession([]byte{}, 60*60)
		data, err := proto.Marshal(session)
		if err != nil {
			return err
		}
		err = tx.Bucket(bucketName).Put([]byte("test1"), data)
		if err != nil {
			return err
		}
		err = tx.Bucket(bucketName).Put([]byte("test2"), data)
		if err != nil {
			return err
		}
		err = tx.Bucket(bucketName).Put([]byte("test3"), data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err.Error())
	}
	options.BatchSize = 3
	go reap(db, options, quitC, doneC)
	time.Sleep(2 * time.Second)
	Quit(quitC, doneC)

}
