package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
)

func (p *PosixFS) GetNetINode(netINodeID solofsapitypes.NetINodeID) (solofsapitypes.NetINodeUintptr, error) {
	return p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
}

func (p *PosixFS) ReleaseNetINode(uNetINode solofsapitypes.NetINodeUintptr) {
	p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
}

func (p *PosixFS) NetINodePWriteWithNetQuery(uNetINode solofsapitypes.NetINodeUintptr,
	netQuery *snettypes.NetQuery, dataLength int, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithNetQuery(uNetINode, netQuery, dataLength, offset)
}

func (p *PosixFS) NetINodePWriteWithMem(uNetINode solofsapitypes.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithMem(uNetINode, data, offset)
}
