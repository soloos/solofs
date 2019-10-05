package metastg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"sync/atomic"
)

func (p *SolodnDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	var (
		backends    snettypes.PeerGroup
		solodnIndex uint32
	)
	solodnIndex = atomic.AddUint32(&p.chooseSolodnIndex, 1)

	backends.Reset()
	p.solodnsForBlockRWMutex.RLock()
	for i := uint32(0); i < 3; i++ {
		solodnIndex = (solodnIndex + uint32(i)) % uint32(len(p.solodnsForBlock))
		backends.Append(p.solodnsForBlock[solodnIndex])
	}
	p.solodnsForBlockRWMutex.RUnlock()
	return backends, nil
}
