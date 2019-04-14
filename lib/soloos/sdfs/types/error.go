package types

import (
	fsapitypes "soloos/common/fsapi/types"

	"golang.org/x/xerrors"
)

var (
	ErrRemoteService      = xerrors.New("remote service error")
	ErrServiceNotExists   = xerrors.New("service not exists")
	ErrObjectExists       = xerrors.New("object exists")
	ErrObjectHasChildren  = xerrors.New("object has children")
	ErrObjectNotExists    = xerrors.New("object not exists")
	ErrObjectNotPrepared  = xerrors.New("object not prepared")
	ErrNetBlockPWrite     = xerrors.New("net block pwrite error")
	ErrNetBlockPRead      = xerrors.New("net block pread error")
	ErrBackendListIsEmpty = xerrors.New("backend list is empty")
	ErrRetryTooManyTimes  = xerrors.New("retry too many times")
	ErrRLockFailed        = xerrors.New("rlock failed")
	ErrLockFailed         = xerrors.New("lock failed")
	ErrInvalidArgs        = xerrors.New("invalid args")
	ErrHasNotPermission   = xerrors.New("has not permission")
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
