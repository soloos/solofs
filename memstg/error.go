package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofstypes"
)

func ErrorToFsStatus(err error) fsapitypes.Status {
	if err == nil {
		return fsapitypes.OK
	}

	switch err.Error() {
	case solofstypes.ErrObjectNotExists.Error():
		return fsapitypes.ENOENT
	case solofstypes.ErrObjectExists.Error():
		return solofstypes.FS_EEXIST
	case solofstypes.ErrRLockFailed.Error():
		return fsapitypes.EAGAIN
	case solofstypes.ErrLockFailed.Error():
		return fsapitypes.EAGAIN
	case solofstypes.ErrObjectHasChildren.Error():
		return solofstypes.FS_ENOTEMPTY
	case solofstypes.ErrHasNotPermission.Error():
		return fsapitypes.EPERM
	}

	return fsapitypes.EIO
}
