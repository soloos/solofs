package solofstypes

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
)

func ErrorToFsStatus(err error) fsapitypes.Status {
	if err == nil {
		return fsapitypes.OK
	}

	switch err.Error() {
	case solofsapitypes.ErrObjectNotExists.Error():
		return fsapitypes.ENOENT
	case solofsapitypes.ErrObjectExists.Error():
		return FS_EEXIST
	case solofsapitypes.ErrRLockFailed.Error():
		return fsapitypes.EAGAIN
	case solofsapitypes.ErrLockFailed.Error():
		return fsapitypes.EAGAIN
	case solofsapitypes.ErrObjectHasChildren.Error():
		return FS_ENOTEMPTY
	case solofsapitypes.ErrHasNotPermission.Error():
		return fsapitypes.EPERM
	}

	return fsapitypes.EIO
}
