package types

import "errors"

var (
	ErrObjectNotExists    = errors.New("object not exists")
	ErrNetBlockPWrite     = errors.New("net block pwrite error")
	ErrNetBlockPRead      = errors.New("net block pread error")
	ErrBackendListIsEmpty = errors.New("backend list is empty")
)
