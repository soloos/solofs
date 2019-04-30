package metastg

import (
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
	"sync"
)

type DataNodeDriver struct {
	offheapDriver *offheap.OffheapDriver
	metaStg       *MetaStg
	snetDriver    *snet.NetDriver

	chooseDataNodeIndex         uint32
	dataNodesForBlockRegistered map[snettypes.PeerID]snettypes.PeerUintptr
	dataNodesForBlock           []snettypes.PeerUintptr
	dataNodesForBlockRWMutex    sync.RWMutex
}

func (p *DataNodeDriver) Init(metaStg *MetaStg, snetDriver *snet.NetDriver) error {
	p.offheapDriver = metaStg.offheapDriver
	p.metaStg = metaStg
	p.snetDriver = snetDriver
	p.dataNodesForBlockRegistered = make(map[snettypes.PeerID]snettypes.PeerUintptr)
	return nil
}

func (p *DataNodeDriver) RegisterDataNode(peerID snettypes.PeerID, addr string) error {
	var (
		uDataNode  snettypes.PeerUintptr
		registered bool
	)

	p.dataNodesForBlockRWMutex.Lock()
	_, registered = p.dataNodesForBlockRegistered[peerID]
	if registered == false {
		uDataNode, _ = p.snetDriver.MustGetPeer(&peerID, addr, types.DefaultSDFSRPCProtocol)
		p.dataNodesForBlockRegistered[peerID] = uDataNode
		p.dataNodesForBlock = append(p.dataNodesForBlock, uDataNode)
	}
	p.dataNodesForBlockRWMutex.Unlock()

	return nil
}

func (p *DataNodeDriver) GetDataNode(peerID snettypes.PeerID) snettypes.PeerUintptr {
	return p.snetDriver.GetPeer(peerID)
}
