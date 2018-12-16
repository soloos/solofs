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
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINode(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsMetaDataInited == false {
		pNetINode.MetaDataInitMutex.Lock()
		if pNetINode.IsMetaDataInited == false {
			err = p.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
			if err == nil {
				pNetINode.IsMetaDataInited = true
			}
		}
		pNetINode.MetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetINode
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINode(netINodeID types.NetINodeID,
	size int64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINode(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsMetaDataInited == false {
		pNetINode.MetaDataInitMutex.Lock()
		if pNetINode.IsMetaDataInited == false {
			err = p.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
			if err == nil {
				pNetINode.IsMetaDataInited = true
			}
		}
		pNetINode.MetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetINode
		return 0, err
	}

	return uNetINode, nil
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
