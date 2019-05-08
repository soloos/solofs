package api

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

// DataNode
type GetDataNode func(peerID snettypes.PeerID) snettypes.PeerUintptr
type ChooseDataNodesForNewNetBlock func(uNetINode types.NetINodeUintptr) (snettypes.PeerGroup, error)

// NetINode
type GetNetINode func(netINodeID types.NetINodeID) (types.NetINodeUintptr, error)
type MustGetNetINode func(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error)
type ReleaseNetINode func(uNetINode types.NetINodeUintptr)
type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode types.NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error
type NetINodeCommitSizeInDB func(uNetINode types.NetINodeUintptr, size uint64) error

// FsINode
type AllocFsINodeID func() types.FsINodeID
type DeleteFsINodeByIDInDB func(fsINodeID types.FsINodeID) error
type ListFsINodeByParentIDFromDB func(parentID types.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(sdfsapitypes.FsINodeMeta) bool,
) error
type UpdateFsINodeInDB func(fsINodeMeta *sdfsapitypes.FsINodeMeta) error
type InsertFsINodeInDB func(fsINodeMeta *sdfsapitypes.FsINodeMeta) error
type FetchFsINodeByIDFromDB func(pFsINodeMeta *sdfsapitypes.FsINodeMeta) error
type FetchFsINodeByNameFromDB func(pFsINodeMeta *sdfsapitypes.FsINodeMeta) error

// FsINodeXAttr
type DeleteFIXAttrInDB func(fsINodeID types.FsINodeID) error
type ReplaceFIXAttrInDB func(fsINodeID types.FsINodeID, xattr types.FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(fsINodeID types.FsINodeID) (types.FsINodeXAttr, error)
