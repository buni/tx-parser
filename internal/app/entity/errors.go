package entity

import "errors"

var (
	ErrBlockHightNotSet = errors.New("block height not set")
	ErrNotFound         = errors.New("not found")
	ErrCurrentBlockNil  = errors.New("current block is nil")
)
