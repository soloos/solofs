package types

import "errors"

var (
	ErrRemoteService      = errors.New("remote service error")
	ErrObjectNotExists    = errors.New("object not exists")
	ErrNetBlockPWrite     = errors.New("net block pwrite error")
	ErrNetBlockPRead      = errors.New("net block pread error")
	ErrBackendListIsEmpty = errors.New("backend list is empty")
	ErrRetryTooManyTimes  = errors.New("retry too many times")
)
