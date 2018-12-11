package api

import (
	"soloos/sdfs/types"
)

func (p *NameNodeClient) PrepareNetINodeMetadata(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int,
) error {
	uNetINode.Ptr().Size = size
	uNetINode.Ptr().NetBlockCap = netBlockCap
	uNetINode.Ptr().MemBlockCap = memBlockCap
	return nil
}
