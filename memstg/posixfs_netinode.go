package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

func (p *PosixFs) GetNetINode(netINodeID solofstypes.NetINodeID) (solofstypes.NetINodeUintptr, error) {
	return p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
}

func (p *PosixFs) ReleaseNetINode(uNetINode solofstypes.NetINodeUintptr) {
	p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
}

func (p *PosixFs) NetINodePWriteWithNetQuery(uNetINode solofstypes.NetINodeUintptr,
	netQuery *snet.NetQuery, dataLength int, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithNetQuery(uNetINode, netQuery, dataLength, offset)
}

func (p *PosixFs) NetINodePWriteWithMem(uNetINode solofstypes.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithMem(uNetINode, data, offset)
}
