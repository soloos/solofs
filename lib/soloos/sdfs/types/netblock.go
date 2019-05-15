package types

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
)

const (
	NetINodeBlockIDSize = sdfsapitypes.NetINodeBlockIDSize
	NetBlockStructSize  = sdfsapitypes.NetBlockStructSize
)

var (
	EncodeNetINodeBlockID = sdfsapitypes.EncodeNetINodeBlockID
)

type NetINodeBlockID = sdfsapitypes.NetINodeBlockID
type NetBlockUintptr = sdfsapitypes.NetBlockUintptr
type NetBlock = sdfsapitypes.NetBlock
