package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

func (p *PosixFs) GetNetINode(netINodeID solofsapitypes.NetINodeID) (solofsapitypes.NetINodeUintptr, error) {
	return p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
}

func (p *PosixFs) ReleaseNetINode(uNetINode solofsapitypes.NetINodeUintptr) {
	p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
}

func (p *PosixFs) NetINodePWriteWithNetQuery(uNetINode solofsapitypes.NetINodeUintptr,
	netQuery *snettypes.NetQuery, dataLength int, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithNetQuery(uNetINode, netQuery, dataLength, offset)
}

func (p *PosixFs) NetINodePWriteWithMem(uNetINode solofsapitypes.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithMem(uNetINode, data, offset)
}
