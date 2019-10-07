package metastg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

func (p *NetBlockDriver) ChooseSolodnsForNewNetBlock(uNetINode solofstypes.NetINodeUintptr) (snet.PeerGroup, error) {
	return p.helper.ChooseSolodnsForNewNetBlock(uNetINode)
}
