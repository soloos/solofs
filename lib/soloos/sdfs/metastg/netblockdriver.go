package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"strings"
)

type NetBlockDriverHelper struct {
	ChooseDataNodesForNewNetBlock sdfsapitypes.ChooseDataNodesForNewNetBlock
	GetDataNode                   sdfsapitypes.GetDataNode
}

type NetBlockDriver struct {
	dbConn *sdbapi.Connection
	helper NetBlockDriverHelper
}

func (p *NetBlockDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getDataNode sdfsapitypes.GetDataNode,
	chooseDataNodesForNewNetBlock sdfsapitypes.ChooseDataNodesForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(getDataNode, chooseDataNodesForNewNetBlock)
	return nil
}

func (p *NetBlockDriver) SetHelper(
	getDataNode sdfsapitypes.GetDataNode,
	chooseDataNodesForNewNetBlock sdfsapitypes.ChooseDataNodesForNewNetBlock,
) {
	p.helper.GetDataNode = getDataNode
	p.helper.ChooseDataNodesForNewNetBlock = chooseDataNodesForNewNetBlock
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netBlockIndex int32) error {
	var (
		pNetBlock           = uNetBlock.Ptr()
		backendPeerIDArrStr string
		peerID              snettypes.PeerID
		err                 error
	)

	err = p.FetchNetBlockFromDB(uNetINode.Ptr(), netBlockIndex, pNetBlock, &backendPeerIDArrStr)
	if err == nil {
		backendPeerIDArr := strings.Split(backendPeerIDArrStr, ",")
		for _, peerIDStr := range backendPeerIDArr {
			copy(peerID[:], peerIDStr)
			pNetBlock.StorDataBackends.Append(peerID)
		}

	} else {
		if err != sdfsapitypes.ErrObjectNotExists {
			goto PREPARE_DONE
		}

		pNetBlock.NetINodeID = uNetINode.Ptr().ID
		pNetBlock.IndexInNetINode = netBlockIndex
		pNetBlock.Len = 0
		pNetBlock.Cap = uNetINode.Ptr().NetBlockCap
		pNetBlock.StorDataBackends, err = p.ChooseDataNodesForNewNetBlock(uNetINode)
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

func (p *NetBlockDriver) PrepareNetBlockSyncDataBackends(uNetBlock sdfsapitypes.NetBlockUintptr,
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

func (p *NetBlockDriver) PrepareNetBlockLocalDataBackend(uNetBlock sdfsapitypes.NetBlockUintptr) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsLocalDataBackendInited.LockContext()
	if pNetBlock.IsLocalDataBackendInited.Load() == sdbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	pNetBlock.IsLocalDataBackendExists = true
	pNetBlock.IsLocalDataBackendInited.Store(sdbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited.UnlockContext()
	return err
}
