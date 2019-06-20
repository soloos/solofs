package memstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
)

func (p *PosixFS) GetNetINode(netINodeID sdfsapitypes.NetINodeID) (sdfsapitypes.NetINodeUintptr, error) {
	return p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
}

func (p *PosixFS) ReleaseNetINode(uNetINode sdfsapitypes.NetINodeUintptr) {
	p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
}

func (p *PosixFS) NetINodePWriteWithNetQuery(uNetINode sdfsapitypes.NetINodeUintptr,
	netQuery *snettypes.NetQuery, dataLength int, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithNetQuery(uNetINode, netQuery, dataLength, offset)
}

func (p *PosixFS) NetINodePWriteWithMem(uNetINode sdfsapitypes.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithMem(uNetINode, data, offset)
}
