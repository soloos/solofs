package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

type DataNodeDriver struct {
	metaStg      *MetaStg
	dataNodePool snettypes.PeerPool
}

func (p *DataNodeDriver) Init(metaStg *MetaStg) error {
	p.metaStg = metaStg
	p.dataNodePool.Init(-1, p.metaStg.offheapDriver)
	return nil
}

func (p *DataNodeDriver) RegisterDataNode(peerID snettypes.PeerID, serveAddr string) (snettypes.PeerUintptr, error) {
	var (
		uDataNode snettypes.PeerUintptr
		exists    bool
	)

	uDataNode, exists = p.dataNodePool.MustGetPeer(peerID)
	if exists == false || uDataNode.Ptr().IsMetaDataInited == false {
		pDataNode := uDataNode.Ptr()
		pDataNode.MetaDataMutex.Lock()
		if pDataNode.IsMetaDataInited == false {
			pDataNode.SetAddress(serveAddr)
			pDataNode.ServiceProtocol = types.DefaultSDFSRPCProtocol
		}
		pDataNode.MetaDataMutex.Unlock()
	}

	return uDataNode, nil
}

func (p *DataNodeDriver) GetDataNode(peerID snettypes.PeerID) (snettypes.PeerUintptr, bool) {
	return p.dataNodePool.GetPeer(peerID)
}
