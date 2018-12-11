package api

import (
	"soloos/sdfs/types"
)

func (p *NameNodeClient) PrepareINodeMetadata(uINode types.INodeUintptr,
	size int64, netBlockCap int, memBlockCap int,
) error {
	uINode.Ptr().Size = size
	uINode.Ptr().NetBlockCap = netBlockCap
	uINode.Ptr().MemBlockCap = memBlockCap
	return nil
}
