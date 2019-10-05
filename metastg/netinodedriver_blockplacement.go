package metastg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

func (p *NetINodeDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	return p.helper.ChooseSolodnsForNewNetBlock(uNetINode)
}
