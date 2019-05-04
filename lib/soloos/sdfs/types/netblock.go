package types

import (
	sdbapitypes "soloos/common/sdbapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
)

const (
	NetINodeBlockIDSize = sdfsapitypes.NetINodeBlockIDSize
	NetBlockStructSize  = sdfsapitypes.NetBlockStructSize
)

type NetINodeBlockID = sdfsapitypes.NetINodeBlockID
type NetBlockUintptr = sdfsapitypes.NetBlockUintptr
type NetBlock = sdfsapitypes.NetBlock
type PtrBindIndex = sdbapitypes.PtrBindIndex

var (
	EncodeNetINodeBlockID = sdfsapitypes.EncodeNetINodeBlockID
)
