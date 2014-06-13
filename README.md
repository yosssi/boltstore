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
// Fetch new store.
store, err := boltstore.New(
	boltstore.Config{
		SessionOptions: sessions.Options{
			MaxAge: 60 * 60 * 24 * 30, // 30days
		},
	},
	[]byte("secret-key"),
)
if err != nil {
	panic(err)
}
defer store.Close()

// Get a session.
sessions, err := store.Get(r, "session-key")
if err != nil {
	panic(err)
}

// Add a value.
session.Values["foo"] = "bar"

// Save.
if err := sessions.Save(r, w); err != nil {
	panic(err)
}

// Delete session.
session.Options.MaxAge = -1
if err := sessions.Save(r, w); err != nil {
	panic(err)
}
```

## Documentation
* [GoDoc](http://godoc.org/github.com/yosssi/boltstore)
