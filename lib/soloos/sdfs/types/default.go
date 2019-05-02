package types

import (
	snettypes "soloos/common/snet/types"
)

const (
	DefaultSDFSRPCNetwork  = "tcp"
	DefaultSDFSRPCProtocol = snettypes.ProtocolSRPC

	DefaultNetBlockCap int = 1024 * 1024 * 8
	DefaultMemBlockCap int = 1024 * 1024 * 2
)
