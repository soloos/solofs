package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
)

type INodeDriver struct {
	offheapDriver  *offheap.OffheapDriver
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *MemBlockDriver
	inodePool      types.INodePool

	nameNodeClient *api.NameNodeClient
}

func (p *INodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver) error {
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.inodePool.Init(-1, p.offheapDriver)
	return nil
}

// MustGetINode get or init a inodeblock
func (p *INodeDriver) MustGetINode(inodeID types.INodeID) (types.INodeUintptr, bool) {
	return p.inodePool.MustGetINode(inodeID)
}

func (p *INodeDriver) InitINode(size int64, netBlockCap int, memBlockCap int) (types.INodeUintptr, error) {
	var (
		inodeID types.INodeID
		uINode  types.INodeUintptr
		exists  bool
		err     error
	)

	util.InitUUID64(&inodeID)
	uINode, exists = p.MustGetINode(inodeID)
	if exists {
		panic("inode should not exists")
	}

	err = p.prepareINodeMetadata(uINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return 0, err
	}

	return uINode, nil
}

func (p *INodeDriver) prepareINodeMetadata(uINode types.INodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var (
		pINode = uINode.Ptr()
		err    error
	)

	pINode.MetaDataMutex.Lock()
	if pINode.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.nameNodeClient.PrepareINodeMetadata(uINode, size, netBlockCap, memBlockCap)
	if err != nil {
		goto PREPARE_DONE
	}

	pINode.IsMetaDataInited = true

PREPARE_DONE:
	pINode.MetaDataMutex.Unlock()
	return err
}
