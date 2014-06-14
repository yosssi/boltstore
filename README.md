# BoltStore

[![GoDoc](https://godoc.org/github.com/yosssi/boltstore?status.png)](https://godoc.org/github.com/yosssi/boltstore)

## About

BoltStore is a session store backend for [gorilla/sessions](https://github.com/gorilla/sessions) using [Bolt](https://github.com/boltdb/bolt).

## Installation

```go
go get github.com/yosssi/boltstore/...
```

## Example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	"github.com/yosssi/boltstore/reaper"
	"github.com/yosssi/boltstore/store"
)

var db *bolt.DB

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	// Fetch a new store.
	str, err := store.New(
		db,
		store.Config{
			SessionOptions: sessions.Options{
				MaxAge: 60 * 60 * 24 * 30, // 30days
			},
		},
		[]byte("secret-key"),
	)

	// Get a session.
	session, err := str.Get(r, "session-key")
	if err != nil {
		panic(err)
	}

	// Add a value.
	session.Values["foo"] = "bar"

	// Save.
	if err := sessions.Save(r, w); err != nil {
		panic(err)
	}

	// Delete the session.
	session.Options.MaxAge = -1
	if err := sessions.Save(r, w); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Hello BoltStore")
}

func main() {
	// Open a Bolt database.
	db, err := bolt.Open("./sessions.db", 0666)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Invoke a reaper which removes expired sessions.
	defer reaper.Quit(reaper.Run(db, reaper.Options{}))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

## Documentation
* [GoDoc](http://godoc.org/github.com/yosssi/boltstore)
