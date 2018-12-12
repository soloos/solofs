package metastg

import (
	"soloos/sdfs/types"
)

type NetINodeDriver struct {
	metaStg      *MetaStg
	netINodePool types.NetINodePool
}

func (p *NetINodeDriver) Init(metaStg *MetaStg) error {
	p.metaStg = metaStg
	p.netINodePool.Init(-1, p.metaStg.offheapDriver)
	return nil
}

func (p *NetINodeDriver) GetNetINode(netINodeID types.NetINodeID) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		exists    bool
		err       error
	)

	uNetINode, exists = p.netINodePool.MustGetNetINode(netINodeID)

	if exists == false || uNetINode.Ptr().IsMetaDataInited == false {
		err = p.prepareNetINodeMetadataOnlyLoadDB(uNetINode)
		if err != nil {
			goto GETINODE_DONE
		}
	}

GETINODE_DONE:
	if err == types.ErrObjectNotExists {
		p.netINodePool.ReleaseNetINode(uNetINode)
	}

	return uNetINode, err
}

func (p *NetINodeDriver) MustGetNetINode(netINodeID types.NetINodeID,
	size int64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		exists    bool
		err       error
	)

	uNetINode, exists = p.netINodePool.MustGetNetINode(netINodeID)

	if exists == false || uNetINode.Ptr().IsMetaDataInited == false {
		err = p.prepareNetINodeMetadata(uNetINode, size, netBlockCap, memBlockCap)
		if err != nil {
			goto GETINODE_DONE
		}
	}

GETINODE_DONE:
	return uNetINode, err
}

func (p *NetINodeDriver) prepareNetINodeMetadataOnlyLoadDB(uNetINode types.NetINodeUintptr) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	pNetINode.MetaDataMutex.Lock()
	if pNetINode.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		goto PREPARE_DONE
	}

	pNetINode.IsMetaDataInited = true

PREPARE_DONE:
	pNetINode.MetaDataMutex.Unlock()
	return err
}

func (p *NetINodeDriver) prepareNetINodeMetadata(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	pNetINode.MetaDataMutex.Lock()
	if pNetINode.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		if err != types.ErrObjectNotExists {
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

	pNetINode.IsMetaDataInited = true

PREPARE_DONE:
	pNetINode.MetaDataMutex.Unlock()
	return err
}
