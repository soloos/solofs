package metastg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

type DataNodeDriver struct {
	*soloosbase.SoloOSEnv
	metaStg *MetaStg

	chooseDataNodeIndex         uint32
	dataNodesForBlockRegistered map[snettypes.PeerID]snettypes.PeerUintptr
	dataNodesForBlock           []snettypes.PeerUintptr
	dataNodesForBlockRWMutex    util.RWMutex
}

func (p *DataNodeDriver) Init(metaStg *MetaStg) error {
	p.SoloOSEnv = metaStg.SoloOSEnv
	p.metaStg = metaStg
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
		uDataNode, _ = p.SNetDriver.MustGetPeer(&peerID, addr, sdfsapitypes.DefaultSDFSRPCProtocol)
		p.dataNodesForBlockRegistered[peerID] = uDataNode
		p.dataNodesForBlock = append(p.dataNodesForBlock, uDataNode)
	}
	p.dataNodesForBlockRWMutex.Unlock()

	return nil
}

func (p *DataNodeDriver) GetDataNode(peerID snettypes.PeerID) snettypes.PeerUintptr {
	return p.SNetDriver.GetPeer(peerID)
}
