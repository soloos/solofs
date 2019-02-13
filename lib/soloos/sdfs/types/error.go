package types

import (
	"errors"
	fsapitypes "soloos/common/fsapi/types"
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
	ErrHasNotPermission   = errors.New("has not permission")
)

func ErrorToFsStatus(err error) fsapitypes.Status {
	switch err {
	case nil:
		return fsapitypes.OK
	case ErrObjectNotExists:
		return fsapitypes.ENOENT
	case ErrObjectExists:
		return FS_EEXIST
	case ErrRLockFailed:
		return fsapitypes.EAGAIN
	case ErrLockFailed:
		return fsapitypes.EAGAIN
	case ErrObjectHasChildren:
		return FS_ENOTEMPTY
	case ErrHasNotPermission:
		return fsapitypes.EPERM
	}

	return fsapitypes.EIO
}
