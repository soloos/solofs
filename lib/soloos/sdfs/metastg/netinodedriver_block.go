package metastg

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr) (snettypes.PeerGroup, error) {
	return p.helper.ChooseDataNodesForNewNetBlock(uNetINode)
}
