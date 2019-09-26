package metastg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
)

func (p *NetBlockDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	return p.helper.ChooseSolodnsForNewNetBlock(uNetINode)
}
