package metastg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"sync"
	"sync/atomic"
)

type DataNodeDriver struct {
	metaStg             *MetaStg
	snetDriver          snet.NetDriver
	chooseDataNodeIndex uint32
	dataNodeRWMutex     sync.RWMutex
	dataNodes           []snettypes.PeerUintptr
}

func (p *DataNodeDriver) Init(metaStg *MetaStg) error {
	p.metaStg = metaStg
	p.snetDriver.Init(p.metaStg.offheapDriver)
	return nil
}

func (p *DataNodeDriver) MustGetDataNode(peerID *snettypes.PeerID, serveAddr string) (snettypes.PeerUintptr, error) {
	var (
		uDataNode snettypes.PeerUintptr
		exists    bool
	)

	uDataNode, exists = p.metaStg.snetDriver.MustGetPeer(peerID, serveAddr, types.DefaultSDFSRPCProtocol)
	if exists == false || uDataNode.Ptr().IsMetaDataInited == false {
		pDataNode := uDataNode.Ptr()
		pDataNode.MetaDataMutex.Lock()

		p.dataNodeRWMutex.Lock()
		p.dataNodes = append(p.dataNodes, uDataNode)
		p.dataNodeRWMutex.Unlock()

		pDataNode.MetaDataMutex.Unlock()
	}

	return uDataNode, nil
}

func (p *DataNodeDriver) GetDataNode(peerID snettypes.PeerID) snettypes.PeerUintptr {
	return p.snetDriver.GetPeer(peerID)
}

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
