package db

import "github.com/prysmaticlabs/prysm/slasher/internal/db/kv"

var _ Database = (*kv.Store)(nil)
