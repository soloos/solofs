package metastg

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

type DataNodeDriver struct {
	*soloosbase.SoloOSEnv
	metaStg *MetaStg

	chooseDataNodeIndex         uint32
	dataNodesForBlockRegistered map[snettypes.PeerID]int64
	dataNodesForBlock           []snettypes.PeerID
	dataNodesForBlockRWMutex    util.RWMutex
}

func (p *DataNodeDriver) Init(metaStg *MetaStg) error {
	p.SoloOSEnv = metaStg.SoloOSEnv
	p.metaStg = metaStg
	p.dataNodesForBlockRegistered = make(map[snettypes.PeerID]int64)
	return nil
}

func (p *DataNodeDriver) DataNodeHeartBeat(peer snettypes.Peer) error {
	var (
		registered bool
	)

	p.dataNodesForBlockRWMutex.Lock()
	_, registered = p.dataNodesForBlockRegistered[peer.ID]
	if registered == false {
		p.SNetDriver.RegisterPeer(peer)
		p.dataNodesForBlockRegistered[peer.ID] = 0
		p.dataNodesForBlock = append(p.dataNodesForBlock, peer.ID)
	}
	p.dataNodesForBlockRWMutex.Unlock()

	return nil
}

func (p *DataNodeDriver) GetDataNode(peerID snettypes.PeerID) (snettypes.Peer, error) {
	return p.SNetDriver.GetPeer(peerID)
}
