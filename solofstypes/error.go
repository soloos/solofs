package solofstypes

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
)

func ErrorToFsStatus(err error) fsapitypes.Status {
	switch err {
	case nil:
		return fsapitypes.OK
	case solofsapitypes.ErrObjectNotExists:
		return fsapitypes.ENOENT
	case solofsapitypes.ErrObjectExists:
		return FS_EEXIST
	case solofsapitypes.ErrRLockFailed:
		return fsapitypes.EAGAIN
	case solofsapitypes.ErrLockFailed:
		return fsapitypes.EAGAIN
	case solofsapitypes.ErrObjectHasChildren:
		return FS_ENOTEMPTY
	case solofsapitypes.ErrHasNotPermission:
		return fsapitypes.EPERM
	}

	return fsapitypes.EIO
}
