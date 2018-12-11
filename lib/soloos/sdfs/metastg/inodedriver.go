package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
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
		err = p.prepareINodeMetadata(uINode)
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

func (p *INodeDriver) ChooseDataNodesForNewNetBlock(uINode types.INodeUintptr,
	backends *snettypes.PeerUintptrArray8) error {
	backends.Reset()
	return nil
}

func (p *INodeDriver) prepareINodeMetadata(uINode types.INodeUintptr) error {
	var (
		pINode = uINode.Ptr()
		err    error
	)

	pINode.MetaDataMutex.Lock()
	if pINode.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.FetchINodeFromDB(pINode)
	goto PREPARE_DONE

	pINode.IsMetaDataInited = true

PREPARE_DONE:
	pINode.MetaDataMutex.Unlock()
	return err
}
