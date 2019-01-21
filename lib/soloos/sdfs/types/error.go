package types

import (
	"errors"

	"github.com/hanwen/go-fuse/fuse"
)

var (
	ErrRemoteService      = errors.New("remote service error")
	ErrServiceNotExists   = errors.New("service not exists")
	ErrObjectExists       = errors.New("object exists")
	ErrObjectHasChildren  = errors.New("object has children")
	ErrObjectNotExists    = errors.New("object not exists")
	ErrObjectNotPrepared  = errors.New("object not prepared")
	ErrNetBlockPWrite     = errors.New("net block pwrite error")
	ErrNetBlockPRead      = errors.New("net block pread error")
	ErrBackendListIsEmpty = errors.New("backend list is empty")
	ErrRetryTooManyTimes  = errors.New("retry too many times")
	ErrRLockFailed        = errors.New("rlock failed")
	ErrLockFailed         = errors.New("lock failed")
	ErrInvalidArgs        = errors.New("invalid args")
)

func ErrorToFuseStatus(err error) fuse.Status {
	switch err {
	case nil:
		return fuse.OK
	case ErrObjectNotExists:
		return fuse.ENOENT
	case ErrObjectExists:
		return FS_EEXIST
	case ErrRLockFailed:
		return fuse.EAGAIN
	case ErrLockFailed:
		return fuse.EAGAIN
	case ErrObjectHasChildren:
		return FS_ENOTEMPTY
	}

	return fuse.EIO
}
