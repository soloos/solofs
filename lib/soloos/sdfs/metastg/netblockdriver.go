package metastg

import (
	"soloos/dbcli"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"strings"
)

type NetBlockDriverHelper struct {
	DBConn                        *dbcli.Connection
	ChooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock
	GetDataNode                   api.GetDataNode
}

type NetBlockDriver struct {
	helper NetBlockDriverHelper
}

func (p *NetBlockDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbcli.Connection,
	getDataNode api.GetDataNode,
	chooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock,
) error {
	p.SetHelper(dbConn, getDataNode, chooseDataNodesForNewNetBlock)
	return nil
}

func (p *NetBlockDriver) SetHelper(dbConn *dbcli.Connection,
	getDataNode api.GetDataNode,
	chooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock,
) {
	p.helper.DBConn = dbConn
	p.helper.GetDataNode = getDataNode
	p.helper.ChooseDataNodesForNewNetBlock = chooseDataNodesForNewNetBlock
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
			uPeer = p.helper.GetDataNode(&peerID)
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
		err = p.helper.ChooseDataNodesForNewNetBlock(uNetINode, &pNetBlock.StorDataBackends)
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
		pNetBlock.IsDBMetaDataInited = true
	}
	return err
}

func (p *NetBlockDriver) PrepareNetBlockSyncDataBackendsWithLock(uNetBlock types.NetBlockUintptr,
	backends snettypes.PeerUintptrArray8,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.MemMetaDataInitMutex.Lock()
	if pNetBlock.IsSyncDataBackendsInited == true {
		goto PREPARE_DONE
	}

	pNetBlock.SyncDataBackends = backends
	pNetBlock.SyncDataPrimaryBackendTransferCount = 0

PREPARE_DONE:
	pNetBlock.IsSyncDataBackendsInited = true
	pNetBlock.MemMetaDataInitMutex.Unlock()
	return err
}

func (p *NetBlockDriver) PrepareNetBlockLocalDataBackendWithLock(uNetBlock types.NetBlockUintptr,
	backend snettypes.PeerUintptr,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.MemMetaDataInitMutex.Lock()
	if pNetBlock.IsLocalDataBackendInited == true {
		goto PREPARE_DONE
	}

	pNetBlock.LocalDataBackend = backend

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited = true
	pNetBlock.MemMetaDataInitMutex.Unlock()
	return err
}
