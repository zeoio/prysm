package kv

import "errors"

var (
	// ErrNotFound for database retrieval.
	ErrNotFound = errors.New("not found")
)
