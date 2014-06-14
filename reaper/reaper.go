package reaper

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/yosssi/boltstore/shared"
)

// Run invokes a reap function as a goroutine.
func Run(db *bolt.DB, options Options) (chan<- struct{}, <-chan struct{}) {
	options.setDefault()
	quitC, doneC := make(chan struct{}), make(chan struct{})
	go reap(db, options, quitC, doneC)
	return quitC, doneC
}

// Quit terminates the reap goroutine.
func Quit(quitC chan<- struct{}, doneC <-chan struct{}) {
	quitC <- struct{}{}
	<-doneC
}

func reap(db *bolt.DB, options Options, quitC <-chan struct{}, doneC chan<- struct{}) {
	var prevKey []byte
	for {
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(options.BucketName)
			if bucket == nil {
				return nil
			}

			c := bucket.Cursor()

			var i int

			for k, v := c.Seek(prevKey); ; k, v = c.Next() {
				// If we hit the end of our sessions then
				// exit and start over next time.
				if k == nil {
					prevKey = nil
					return nil
				}

				i++

				session, err := shared.Session(v)
				if err != nil {
					return err
				}

				if shared.Expired(session) {
					err := db.Update(func(txu *bolt.Tx) error {
						return txu.Bucket(options.BucketName).Delete(k)
					})
					if err != nil {
						return err
					}
				}

				if options.BatchSize == i {
					copy(prevKey, k)
					return nil
				}
			}
		})

		if err != nil {
			log.Println(err.Error())
		}

		// Check if a quit signal is sent.
		select {
		case <-quitC:
			doneC <- struct{}{}
			return
		default:
		}

		time.Sleep(options.CheckInterval)
	}
}
