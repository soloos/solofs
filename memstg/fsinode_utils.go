package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func FsModeToFsINodeType(mode uint32) int {
	if mode&fsapi.S_IFDIR != 0 {
		return solofstypes.FSINODE_TYPE_DIR
	}
	if mode&fsapi.S_IFREG != 0 {
		return solofstypes.FSINODE_TYPE_FILE
	}
	if mode&fsapi.S_IFIFO != 0 {
		return solofstypes.FSINODE_TYPE_FIFO
	}
	if mode&fsapi.S_IFLNK != 0 {
		return solofstypes.FSINODE_TYPE_SOFT_LINK
	}
	return solofstypes.FSINODE_TYPE_UNKOWN
}

func FsINodeTypeToFsType(fsINodeType int) int {
	switch fsINodeType {
	case solofstypes.FSINODE_TYPE_DIR:
		return fsapi.S_IFDIR
	case solofstypes.FSINODE_TYPE_FILE:
		return fsapi.S_IFREG
	case solofstypes.FSINODE_TYPE_FIFO:
		return fsapi.S_IFIFO
	case solofstypes.FSINODE_TYPE_SOFT_LINK:
		return fsapi.S_IFLNK
	}
	return fsapi.S_IFREG
}
