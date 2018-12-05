package types

import "errors"

var (
	ErrNetBlockPWrite     = errors.New("net block pwrite error")
	ErrBackendListIsEmpty = errors.New("backend list is empty")
)
