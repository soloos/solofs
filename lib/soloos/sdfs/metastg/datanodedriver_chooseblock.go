package metastg

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
	"sync/atomic"
)

func (p *DataNodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr) (snettypes.PeerUintptrArray8, error) {
	var (
		backends      snettypes.PeerUintptrArray8
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
	p.dataNodesForBlockRWMutex.RLock()
	return backends, nil
}
