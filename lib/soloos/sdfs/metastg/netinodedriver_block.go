package metastg

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr) (snettypes.PeerUintptrArray8, error) {
	return p.helper.ChooseDataNodesForNewNetBlock(uNetINode)
}
