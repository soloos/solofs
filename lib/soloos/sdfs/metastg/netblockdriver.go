package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
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

func (p *NetBlockDriver) GetNetBlock(uINode types.INodeUintptr, netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		exists    bool
		err       error
	)

	uNetBlock, exists = p.netBlockPool.MustGetNetBlock(uINode, netBlockIndex)

	if exists == false || uNetBlock.Ptr().IsMetaDataInited == false {
		uNetBlock.Ptr().MetaDataMutex.Lock()
		if uNetBlock.Ptr().IsMetaDataInited == false {
			err = p.FetchNetBlockFromDB(uNetBlock.Ptr())
		}
		uNetBlock.Ptr().MetaDataMutex.Unlock()
		if err != nil {
			goto GETINODE_DONE
		}
	}

GETINODE_DONE:
	if err == types.ErrObjectNotExists {
		p.netBlockPool.ReleaseNetBlock(uINode, netBlockIndex, uNetBlock)
	}

	return uNetBlock, err
}

func (p *NetBlockDriver) AllocNetBlock(uINode types.INodeUintptr, netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		pNetBlock *types.NetBlock
		err       error
	)

	uNetBlock = p.netBlockPool.AllocRawNetBlock()
	pNetBlock = uNetBlock.Ptr()
	util.InitUUID64(&pNetBlock.ID)
	// pNetBlock.Size = 0
	// pNetBlock.NetBlockCap = netBlockCap
	// pNetBlock.MemBlockCap = memBlockCap
	// pNetBlock.IsMetaDataInited = true

	err = p.StoreNetBlockInDB(uINode.Ptr(), pNetBlock)
	if err != nil {
		goto ALLOCINODE_DONE
	}

	p.netBlockPool.SetNetBlock(uINode, netBlockIndex, uNetBlock)

ALLOCINODE_DONE:
	if err != nil {
		p.netBlockPool.ReleaseRawNetBlock(uNetBlock)
		return 0, err
	}
	return uNetBlock, nil
}
