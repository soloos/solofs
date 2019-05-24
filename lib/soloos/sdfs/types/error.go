package types

import (
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
)

func ErrorToFsStatus(err error) fsapitypes.Status {
	switch err {
	case nil:
		return fsapitypes.OK
	case sdfsapitypes.ErrObjectNotExists:
		return fsapitypes.ENOENT
	case sdfsapitypes.ErrObjectExists:
		return FS_EEXIST
	case sdfsapitypes.ErrRLockFailed:
		return fsapitypes.EAGAIN
	case sdfsapitypes.ErrLockFailed:
		return fsapitypes.EAGAIN
	case sdfsapitypes.ErrObjectHasChildren:
		return FS_ENOTEMPTY
	case sdfsapitypes.ErrHasNotPermission:
		return fsapitypes.EPERM
	}

	return fsapitypes.EIO
}
