package metastg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

func (p *NetBlockDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	return p.helper.ChooseSolodnsForNewNetBlock(uNetINode)
}
