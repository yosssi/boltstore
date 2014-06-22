# BoltStore - Session store using Bolt

[![wercker status](https://app.wercker.com/status/752959ce0f923476671e49fb9b76ebe0/m "wercker status")](https://app.wercker.com/project/bykey/752959ce0f923476671e49fb9b76ebe0)
[![Coverage Status](https://coveralls.io/repos/yosssi/boltstore/badge.png?branch=HEAD)](https://coveralls.io/r/yosssi/boltstore)
[![GoDoc](https://godoc.org/github.com/yosssi/boltstore?status.png)](https://godoc.org/github.com/yosssi/boltstore)

## Overview

BoltStore is a session store using [Bolt](https://github.com/boltdb/bolt) which is a pure Go key/value store. You can store session data in Bolt by using this store. This store implements the [gorilla/sessions](https://github.com/gorilla/sessions) package's [Store](http://godoc.org/github.com/gorilla/sessions#Store) interface. BoltStore's APIs and examples can be seen on its [GoDoc](http://godoc.org/github.com/yosssi/boltstore) page.

## Installation

```go
go get github.com/yosssi/boltstore/...
```

## Example

Here is a simple example using BoltStore. You can see other examples on the BoltStore's [GoDoc](http://godoc.org/github.com/yosssi/boltstore) page.

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
	// Create a store.
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

	fmt.Fprintf(w, "Hello BoltStore")
}

func main() {
	var err error
	// Open a Bolt database.
	db, err = bolt.Open("./sessions.db", 0666, nil)
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

## Benchmarks

```sh
BenchmarkNew	    			5000	    316700 ns/op	   19003 B/op	      35 allocs/op
BenchmarkStore_Get			20000000	       104 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore_New			10000000	       294 ns/op	     130 B/op	       2 allocs/op
BenchmarkStore_Save	    		5000	    488683 ns/op	   65484 B/op	     136 allocs/op
BenchmarkStore_Save_delete	    5000	    476563 ns/op	   59576 B/op	      76 allocs/op
```

## Documentation
* [GoDoc](http://godoc.org/github.com/yosssi/boltstore)
