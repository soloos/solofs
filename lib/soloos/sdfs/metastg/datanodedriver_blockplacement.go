package metastg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"sync/atomic"
)

func (p *DataNodeDriver) ChooseDataNodesForNewNetBlock(uNetINode sdfsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	var (
		backends      snettypes.PeerGroup
		dataNodeIndex uint32
		uDataNode     snettypes.PeerUintptr
	)
	dataNodeIndex = atomic.AddUint32(&p.chooseDataNodeIndex, 1)

	backends.Reset()
	p.dataNodesForBlockRWMutex.RLock()
	for i := uint32(0); i < 3; i++ {
		dataNodeIndex = (dataNodeIndex + uint32(i)) % uint32(len(p.dataNodesForBlock))
		uDataNode = p.dataNodesForBlock[dataNodeIndex]
		backends.Append(uDataNode)
	}
	p.dataNodesForBlockRWMutex.RUnlock()
	return backends, nil
}
