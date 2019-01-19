package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *NetINodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr,
	backends *snettypes.PeerUintptrArray8) error {
	return p.helper.ChooseDataNodesForNewNetBlock(uNetINode, backends)
}
