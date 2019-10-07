package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func ErrorToFsStatus(err error) fsapi.Status {
	if err == nil {
		return fsapi.OK
	}

	switch err.Error() {
	case solofstypes.ErrObjectNotExists.Error():
		return fsapi.ENOENT
	case solofstypes.ErrObjectExists.Error():
		return solofstypes.FS_EEXIST
	case solofstypes.ErrRLockFailed.Error():
		return fsapi.EAGAIN
	case solofstypes.ErrLockFailed.Error():
		return fsapi.EAGAIN
	case solofstypes.ErrObjectHasChildren.Error():
		return solofstypes.FS_ENOTEMPTY
	case solofstypes.ErrHasNotPermission.Error():
		return fsapi.EPERM
	}

	return fsapi.EIO
}
