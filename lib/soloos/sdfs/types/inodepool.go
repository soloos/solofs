package types

import (
	"soloos/util/offheap"
	"sync"
)

type INodePool struct {
	offheapDriver   *offheap.OffheapDriver
	inodeObjectPool offheap.RawObjectPool
	poolRWMutex     sync.RWMutex
	pool            map[INodeID]INodeUintptr
}

func (p *INodePool) Init(rawChunksLimit int32,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.inodeObjectPool,
		int(INodeStructSize), rawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[INodeID]INodeUintptr)

	return nil
}

func (p *INodePool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *INodePool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetINode get or init a inodeblock
func (p *INodePool) MustGetINode(inodeID INodeID) (INodeUintptr, bool) {
	var (
		ret    INodeUintptr
		exists bool
	)

	p.poolRWMutex.RLock()
	ret, exists = p.pool[inodeID]
	p.poolRWMutex.RUnlock()
	if exists {
		return ret, true
	}

	p.poolRWMutex.Lock()
	ret, exists = p.pool[inodeID]
	if exists {
		goto GET_DONE
	}

	ret = INodeUintptr(p.inodeObjectPool.AllocRawObject())
	p.pool[inodeID] = ret

GET_DONE:
	p.poolRWMutex.Unlock()
	return ret, exists
}

func (p *INodePool) ReleaseINode(uINode INodeUintptr) {
	p.poolRWMutex.Lock()
	delete(p.pool, uINode.Ptr().ID)
	p.inodeObjectPool.ReleaseRawObject(uintptr(uINode))
	p.poolRWMutex.Unlock()
}

func (p *INodePool) SetINode(uINode INodeUintptr) {
	p.poolRWMutex.Lock()
	p.pool[uINode.Ptr().ID] = uINode
	p.poolRWMutex.Unlock()
}

func (p *INodePool) AllocRawINode() INodeUintptr {
	return INodeUintptr(p.inodeObjectPool.AllocRawObject())
}

func (p *INodePool) ReleaseRawINode(uINode INodeUintptr) {
	p.inodeObjectPool.ReleaseRawObject(uintptr(uINode))
}
