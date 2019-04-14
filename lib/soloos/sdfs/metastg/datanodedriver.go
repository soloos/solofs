package metastg

import (
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
	"sync"
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
