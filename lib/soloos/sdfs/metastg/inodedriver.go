package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
)

type INodeDriver struct {
	metaStg   *MetaStg
	inodePool types.INodePool
}

func (p *INodeDriver) Init(metaStg *MetaStg) error {
	p.metaStg = metaStg
	p.inodePool.Init(-1, p.metaStg.offheapDriver)
	return nil
}

func (p *INodeDriver) GetINode(inodeID types.INodeID) (types.INodeUintptr, error) {
	var (
		uINode types.INodeUintptr
		exists bool
		err    error
	)

	uINode, exists = p.inodePool.MustGetINode(inodeID)

	if exists == false || uINode.Ptr().IsMetaDataInited == false {
		uINode.Ptr().MetaDataMutex.Lock()
		if uINode.Ptr().IsMetaDataInited == false {
			err = p.FetchINodeFromDB(uINode.Ptr())
		}
		uINode.Ptr().MetaDataMutex.Unlock()
		if err != nil {
			goto GETINODE_DONE
		}
	}

GETINODE_DONE:
	if err == types.ErrObjectNotExists {
		p.inodePool.ReleaseINode(uINode)
	}

	return uINode, err
}

func (p *INodeDriver) AllocINode(netBlockCap, memBlockCap int) (types.INodeUintptr, error) {
	var (
		uINode types.INodeUintptr
		pINode *types.INode
		err    error
	)

	uINode = p.inodePool.AllocRawINode()
	pINode = uINode.Ptr()
	util.InitUUID64(&pINode.ID)
	pINode.Size = 0
	pINode.NetBlockCap = netBlockCap
	pINode.MemBlockCap = memBlockCap
	pINode.IsMetaDataInited = true

	err = p.StoreINodeInDB(pINode)
	if err != nil {
		goto ALLOCINODE_DONE
	}

	p.inodePool.SetINode(uINode)

ALLOCINODE_DONE:
	if err != nil {
		p.inodePool.ReleaseRawINode(uINode)
		return 0, err
	}
	return uINode, nil
}
