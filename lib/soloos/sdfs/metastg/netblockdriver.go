package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"
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

func (p *NetBlockDriver) MustGetNetBlock(uNetINode types.NetINodeUintptr, netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		exists    bool
		err       error
	)

	uNetBlock, exists = p.netBlockPool.MustGetNetBlock(uNetINode, netBlockIndex)

	if exists == false || uNetBlock.Ptr().IsMetaDataInited == false {
		err = p.prepareNetBlockMetadata(uNetINode, netBlockIndex, uNetBlock)
		if err != nil {
			goto GETINODE_DONE
		}
	}

GETINODE_DONE:

	return uNetBlock, err
}

func (p *NetBlockDriver) prepareNetBlockMetadata(uNetINode types.NetINodeUintptr, netBlockIndex int,
	uNetBlock types.NetBlockUintptr) error {
	var (
		pNetBlock           = uNetBlock.Ptr()
		backendPeerIDArrStr string
		peerID              snettypes.PeerID
		err                 error
	)

	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.FetchNetBlockFromDB(pNetBlock, &backendPeerIDArrStr)
	if err == nil {
		// TODO backendPeerIDArrStr split by ','
		backendPeerIDArr := strings.Split(backendPeerIDArrStr, ",")
		for _, peerIDStr := range backendPeerIDArr {
			copy(peerID[:], peerIDStr)
		}

	} else {
		if err != types.ErrObjectNotExists {
			goto PREPARE_DONE
		}

		util.InitUUID64(&pNetBlock.ID)
		pNetBlock.IndexInInode = netBlockIndex
		pNetBlock.Len = 0
		pNetBlock.Cap = uNetINode.Ptr().NetBlockCap
		err = p.metaStg.NetINodeDriver.ChooseDataNodesForNewNetBlock(uNetINode, &pNetBlock.DataNodes)
		if err != nil {
			goto PREPARE_DONE
		}

		err = p.StoreNetBlockInDB(uNetINode.Ptr(), pNetBlock)
		if err != nil {
			goto PREPARE_DONE
		}
	}

	pNetBlock.IsMetaDataInited = true

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
	return err
}
