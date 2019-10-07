package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofstypes"
)

func FsModeToFsINodeType(mode uint32) int {
	if mode&fsapitypes.S_IFDIR != 0 {
		return solofstypes.FSINODE_TYPE_DIR
	}
	if mode&fsapitypes.S_IFREG != 0 {
		return solofstypes.FSINODE_TYPE_FILE
	}
	if mode&fsapitypes.S_IFIFO != 0 {
		return solofstypes.FSINODE_TYPE_FIFO
	}
	if mode&fsapitypes.S_IFLNK != 0 {
		return solofstypes.FSINODE_TYPE_SOFT_LINK
	}
	return solofstypes.FSINODE_TYPE_UNKOWN
}

func FsINodeTypeToFsType(fsINodeType int) int {
	switch fsINodeType {
	case solofstypes.FSINODE_TYPE_DIR:
		return fsapitypes.S_IFDIR
	case solofstypes.FSINODE_TYPE_FILE:
		return fsapitypes.S_IFREG
	case solofstypes.FSINODE_TYPE_FIFO:
		return fsapitypes.S_IFIFO
	case solofstypes.FSINODE_TYPE_SOFT_LINK:
		return fsapitypes.S_IFLNK
	}
	return fsapitypes.S_IFREG
}
