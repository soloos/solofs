package types

import "errors"

var (
	ErrRemoteService      = errors.New("remote service error")
	ErrServiceNotExists   = errors.New("service not exists")
	ErrObjectNotExists    = errors.New("object not exists")
	ErrObjectNotPrepared  = errors.New("object not prepared")
	ErrNetBlockPWrite     = errors.New("net block pwrite error")
	ErrNetBlockPRead      = errors.New("net block pread error")
	ErrBackendListIsEmpty = errors.New("backend list is empty")
	ErrRetryTooManyTimes  = errors.New("retry too many times")
)
