package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *NetINodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr) (snettypes.PeerUintptrArray8, error) {
	return p.helper.ChooseDataNodesForNewNetBlock(uNetINode)
}
