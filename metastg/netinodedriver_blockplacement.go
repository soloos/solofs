package metastg

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
)

func (p *NetINodeDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snet.PeerGroup, error) {
	return p.helper.ChooseSolodnsForNewNetBlock(uNetINode)
}
