package metastg

import (
	"soloos/common/snet"
	"soloos/common/solodbapi"
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"strings"
)

type NetBlockDriverHelper struct {
	ChooseSolodnsForNewNetBlock solofsapitypes.ChooseSolodnsForNewNetBlock
	GetSolodn                   solofsapitypes.GetSolodn
}

type NetBlockDriver struct {
	dbConn *solodbapi.Connection
	helper NetBlockDriverHelper
}

func (p *NetBlockDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
	getSolodn solofsapitypes.GetSolodn,
	chooseSolodnsForNewNetBlock solofsapitypes.ChooseSolodnsForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(getSolodn, chooseSolodnsForNewNetBlock)
	return nil
}

func (p *NetBlockDriver) SetHelper(
	getSolodn solofsapitypes.GetSolodn,
	chooseSolodnsForNewNetBlock solofsapitypes.ChooseSolodnsForNewNetBlock,
) {
	p.helper.GetSolodn = getSolodn
	p.helper.ChooseSolodnsForNewNetBlock = chooseSolodnsForNewNetBlock
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netBlockIndex int32) error {
	var (
		pNetBlock           = uNetBlock.Ptr()
		backendPeerIDArrStr string
		peerID              snet.PeerID
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
		if err.Error() != solofsapitypes.ErrObjectNotExists.Error() {
			goto PREPARE_DONE
		}

		pNetBlock.NetINodeID = uNetINode.Ptr().ID
		pNetBlock.IndexInNetINode = netBlockIndex
		pNetBlock.Len = 0
		pNetBlock.Cap = uNetINode.Ptr().NetBlockCap
		pNetBlock.StorDataBackends, err = p.ChooseSolodnsForNewNetBlock(uNetINode)
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
		pNetBlock.IsDBMetaDataInited.Store(solodbapitypes.MetaDataStateInited)
	}
	return err
}
