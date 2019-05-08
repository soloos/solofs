package metastg

import (
	"soloos/common/sdbapi"
	sdbapitypes "soloos/common/sdbapi/types"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"strings"
)

type NetBlockDriverHelper struct {
	ChooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock
	GetDataNode                   api.GetDataNode
}

type NetBlockDriver struct {
	dbConn *sdbapi.Connection
	helper NetBlockDriverHelper
}

func (p *NetBlockDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getDataNode api.GetDataNode,
	chooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(getDataNode, chooseDataNodesForNewNetBlock)
	return nil
}

func (p *NetBlockDriver) SetHelper(
	getDataNode api.GetDataNode,
	chooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock,
) {
	p.helper.GetDataNode = getDataNode
	p.helper.ChooseDataNodesForNewNetBlock = chooseDataNodesForNewNetBlock
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netBlockIndex int32) error {
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
			uPeer = p.helper.GetDataNode(peerID)
			if uPeer == 0 {
				err = types.ErrObjectNotExists
				goto PREPARE_DONE
			}
			pNetBlock.StorDataBackends.Append(uPeer)
		}

	} else {
		if err != types.ErrObjectNotExists {
			goto PREPARE_DONE
		}

		pNetBlock.NetINodeID = uNetINode.Ptr().ID
		pNetBlock.IndexInNetINode = netBlockIndex
		pNetBlock.Len = 0
		pNetBlock.Cap = uNetINode.Ptr().NetBlockCap
		pNetBlock.StorDataBackends, err = p.helper.ChooseDataNodesForNewNetBlock(uNetINode)
		if err != nil {
			goto PREPARE_DONE
		}

		err = p.StoreNetBlockInDB(uNetINode.Ptr(), pNetBlock)
		if err != nil {
			goto PREPARE_DONE
		}
	}

PREPARE_DONE:
	if err == nil {
		pNetBlock.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	}
	return err
}

func (p *NetBlockDriver) PrepareNetBlockSyncDataBackendsWithLock(uNetBlock types.NetBlockUintptr,
	backends snettypes.PeerGroup,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsSyncDataBackendsInited.LockContext()
	if pNetBlock.IsSyncDataBackendsInited.Load() == sdbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	pNetBlock.SyncDataBackends = backends
	pNetBlock.SyncDataPrimaryBackendTransferCount = 0
	pNetBlock.IsSyncDataBackendsInited.Store(sdbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsSyncDataBackendsInited.UnlockContext()
	return err
}

func (p *NetBlockDriver) PrepareNetBlockLocalDataBackendWithLock(uNetBlock types.NetBlockUintptr,
	backend snettypes.PeerUintptr,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsLocalDataBackendInited.LockContext()
	if pNetBlock.IsLocalDataBackendInited.Load() == sdbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	pNetBlock.LocalDataBackend = backend
	pNetBlock.IsLocalDataBackendInited.Store(sdbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited.UnlockContext()
	return err
}
