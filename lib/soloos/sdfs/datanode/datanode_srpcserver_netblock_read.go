package datanode

import (
	snettypes "soloos/snet/types"
)

func (p *DataNodeSRPCServer) NetBlockPRead(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	return nil
}
