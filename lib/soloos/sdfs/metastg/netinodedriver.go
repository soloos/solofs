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

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode types.NetINodeUintptr) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		if err != types.ErrObjectNotExists {
			return err
		}

		pNetINode.Size = size
		pNetINode.NetBlockCap = netBlockCap
		pNetINode.MemBlockCap = memBlockCap
		err = p.StoreNetINodeInDB(pNetINode)
		if err != nil {
			return err
		}
	}

	return nil
}
