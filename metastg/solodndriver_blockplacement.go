package metastg

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
	"sync/atomic"
)

func (p *SolodnDriver) ChooseSolodnsForNewNetBlock(uNetINode solofsapitypes.NetINodeUintptr) (snet.PeerGroup, error) {
	var (
		backends    snet.PeerGroup
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
