package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"strings"
)

type NetBlockDriver struct {
	metaStg      *MetaStg
	netBlockPool types.NetBlockPool
}

func (p *NetBlockDriver) Init(metaStg *MetaStg) error {
	p.metaStg = metaStg
	p.netBlockPool.Init(-1, p.metaStg.offheapDriver)
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netBlockIndex int) error {
	var (
		pNetBlock           = uNetBlock.Ptr()
		backendPeerIDArrStr string
		peerID              snettypes.PeerID
		uPeer               snettypes.PeerUintptr
		err                 error
	)

	err = p.FetchNetBlockFromDB(uNetINode.Ptr(), netBlockIndex, pNetBlock, &backendPeerIDArrStr)
	if err == nil {
		backendPeerIDArr := strings.Split(backendPeerIDArrStr, ",")
		for _, peerIDStr := range backendPeerIDArr {
			copy(peerID[:], peerIDStr)
			uPeer = p.metaStg.DataNodeDriver.GetDataNode(&peerID)
			if uPeer == 0 {
				return types.ErrObjectNotExists
			}
			pNetBlock.StorDataBackends.Append(uPeer)
		}

	} else {
		if err != types.ErrObjectNotExists {
			return err
		}

		pNetBlock.NetINodeID = uNetINode.Ptr().ID
		pNetBlock.IndexInNetINode = netBlockIndex
		pNetBlock.Len = 0
		pNetBlock.Cap = uNetINode.Ptr().NetBlockCap
		err = p.metaStg.NetINodeDriver.ChooseDataNodesForNewNetBlock(uNetINode, &pNetBlock.StorDataBackends)
		if err != nil {
			return err
		}

		err = p.StoreNetBlockInDB(uNetINode.Ptr(), pNetBlock)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *NetBlockDriver) PrepareNetBlockSyncDataBackendsWithLock(uNetBlock types.NetBlockUintptr,
	backends snettypes.PeerUintptrArray8,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.SyncDataBackendsInitMutex.Lock()
	if pNetBlock.IsSyncDataBackendsInited == true {
		goto PREPARE_DONE
	}

	pNetBlock.SyncDataBackends = backends
	pNetBlock.SyncDataPrimaryBackendTransferCount = 0

	pNetBlock.IsSyncDataBackendsInited = true

PREPARE_DONE:
	pNetBlock.SyncDataBackendsInitMutex.Unlock()
	return err
}
