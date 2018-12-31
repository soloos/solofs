package api

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

type GetNetINode func(netINodeID types.NetINodeID) (types.NetINodeUintptr, error)
type MustGetNetINode func(netINodeID types.NetINodeID,
	size int64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error)
type FetchAndUpdateMaxID func(key string, delta int64) (int64, error)
type GetDataNode func(peerID *snettypes.PeerID) snettypes.PeerUintptr
type ChooseOneDataNode func() snettypes.PeerUintptr
type ChooseDataNodesForNewNetBlock func(uNetINode types.NetINodeUintptr,
	backends *snettypes.PeerUintptrArray8) error

type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode types.NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error

type NetINodeDriverHelper struct {
	NameNodeClient                    *NameNodeClient
	PrepareNetINodeMetaDataOnlyLoadDB PrepareNetINodeMetaDataOnlyLoadDB
	PrepareNetINodeMetaDataWithStorDB PrepareNetINodeMetaDataWithStorDB
}
