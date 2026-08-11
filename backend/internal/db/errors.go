package db

import "errors"

// ErrNotFound is returned by the Get* methods when no item exists for the
// given key.
var ErrNotFound = errors.New("db: not found")
