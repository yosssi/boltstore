# BoltStore - Session store using Bolt

[![wercker status](https://app.wercker.com/status/752959ce0f923476671e49fb9b76ebe0/m "wercker status")](https://app.wercker.com/project/bykey/752959ce0f923476671e49fb9b76ebe0)
[![Coverage Status](https://coveralls.io/repos/yosssi/boltstore/badge.png?branch=HEAD)](https://coveralls.io/r/yosssi/boltstore)
[![GoDoc](https://godoc.org/github.com/yosssi/boltstore?status.png)](https://godoc.org/github.com/yosssi/boltstore)

## About

BoltStore is a session store using [Bolt](https://github.com/boltdb/bolt). This store implements the [gorilla/sessions](https://github.com/gorilla/sessions) package's [Store](http://godoc.org/github.com/gorilla/sessions#Store) interface.

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
	// Fetch a new store.
	str, err := store.New(db, store.Config{}, []byte("secret-key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Get a session.
	session, err := str.Get(r, "session-key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Add a value on the session.
	session.Values["foo"] = "bar"

	// Save the session.
	if err := sessions.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Delete the session.
	session.Options.MaxAge = -1
	if err := sessions.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Hello BoltStore")
}

func main() {
	var err error
	// Open a Bolt database.
	db, err = bolt.Open("./sessions.db", 0666)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Invoke a reaper which checks and removes expired sessions periodically.
	defer reaper.Quit(reaper.Run(db, reaper.Options{}))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

## Documentation
* [GoDoc](http://godoc.org/github.com/yosssi/boltstore)
