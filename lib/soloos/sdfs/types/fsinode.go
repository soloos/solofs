package types

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
)

type FsINodeID = sdfsapitypes.FsINodeID
type DirTreeTime = sdfsapitypes.DirTreeTime
type DirTreeTimeNsec = sdfsapitypes.DirTreeTimeNsec

const (
	ZombieFsINodeParentID = sdfsapitypes.ZombieFsINodeParentID
	RootFsINodeParentID   = sdfsapitypes.RootFsINodeParentID
	RootFsINodeID         = sdfsapitypes.RootFsINodeID
	FsINodeStructSize     = sdfsapitypes.FsINodeStructSize
	MaxFsINodeID          = sdfsapitypes.MaxFsINodeID
)

type FsINodeUintptr = sdfsapitypes.FsINodeUintptr

type FsINode = sdfsapitypes.FsINode

func FsModeToFsINodeType(mode uint32) int {
	if mode&fsapitypes.S_IFDIR != 0 {
		return FSINODE_TYPE_DIR
	}
	if mode&fsapitypes.S_IFREG != 0 {
		return FSINODE_TYPE_FILE
	}
	if mode&fsapitypes.S_IFIFO != 0 {
		return FSINODE_TYPE_FIFO
	}
	if mode&fsapitypes.S_IFLNK != 0 {
		return FSINODE_TYPE_SOFT_LINK
	}
	return FSINODE_TYPE_UNKOWN
}

func FsINodeTypeToFsType(fsINodeType int) int {
	switch fsINodeType {
	case FSINODE_TYPE_DIR:
		return fsapitypes.S_IFDIR
	case FSINODE_TYPE_FILE:
		return fsapitypes.S_IFREG
	case FSINODE_TYPE_FIFO:
		return fsapitypes.S_IFIFO
	case FSINODE_TYPE_SOFT_LINK:
		return fsapitypes.S_IFLNK
	}
	return fsapitypes.S_IFREG
}
