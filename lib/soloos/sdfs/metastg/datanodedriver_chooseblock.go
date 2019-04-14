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
	)
	dataNodeIndex = atomic.AddUint32(&p.chooseDataNodeIndex, 1)

	backends.Reset()
	for i := uint32(0); i < 3; i++ {
		backends.Append(p.dataNodes[int((dataNodeIndex+i)%uint32(len(p.dataNodes)))])
	}
	return backends, nil
}
