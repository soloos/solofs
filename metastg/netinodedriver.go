package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
)

type NetINodeDriverHelper struct {
	ChooseSolodnsForNewNetBlock solofstypes.ChooseSolodnsForNewNetBlock
}

type NetINodeDriver struct {
	*soloosbase.SoloosEnv
	dbConn *solodbapi.Connection
	helper NetINodeDriverHelper
}

func (p *NetINodeDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
	chooseOneSolodn solofstypes.ChooseSolodnsForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(chooseOneSolodn)
	return nil
}

func (p *NetINodeDriver) SetHelper(
	chooseOneSolodn solofstypes.ChooseSolodnsForNewNetBlock,
) {
	p.helper.ChooseSolodnsForNewNetBlock = chooseOneSolodn
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode solofstypes.NetINodeUintptr) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		goto PREPARE_DONE
	}

PREPARE_DONE:
	if err == nil {
		pNetINode.IsDBMetaDataInited.Store(solodbtypes.MetaDataStateInited)
	}
	return err
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode solofstypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		if err.Error() != solofstypes.ErrObjectNotExists.Error() {
			goto PREPARE_DONE
		}

		pNetINode.Size = size
		pNetINode.NetBlockCap = netBlockCap
		pNetINode.MemBlockCap = memBlockCap
		err = p.StoreNetINodeInDB(pNetINode)
		if err != nil {
			goto PREPARE_DONE
		}
	}

PREPARE_DONE:
	if err == nil {
		pNetINode.IsDBMetaDataInited.Store(solodbtypes.MetaDataStateInited)
	}
	return err
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode solofstypes.NetINodeUintptr, size uint64) error {
	return p.NetINodeCommitSizeInDB(uNetINode, size)
}
