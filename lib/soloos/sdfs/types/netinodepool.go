package types

import (
	"soloos/util"
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
		uINode INodeUintptr
		exists bool
	)

	p.poolRWMutex.RLock()
	uINode, exists = p.pool[inodeID]
	p.poolRWMutex.RUnlock()
	if exists {
		return uINode, true
	}

	p.poolRWMutex.Lock()
	uINode, exists = p.pool[inodeID]
	if exists {
		goto GET_DONE
	}

	uINode = INodeUintptr(p.inodeObjectPool.AllocRawObject())
	uINode.Ptr().ID = inodeID
	p.pool[inodeID] = uINode

GET_DONE:
	p.poolRWMutex.Unlock()
	return uINode, exists
}

func (p *INodePool) ReleaseINode(uINode INodeUintptr) {
	var exists bool
	p.poolRWMutex.Lock()
	_, exists = p.pool[uINode.Ptr().ID]
	if exists {
		delete(p.pool, uINode.Ptr().ID)
		p.inodeObjectPool.ReleaseRawObject(uintptr(uINode))
	}
	p.poolRWMutex.Unlock()
}

func (p *INodePool) SaveRawINode(uINode INodeUintptr) {
	p.poolRWMutex.Lock()
	p.pool[uINode.Ptr().ID] = uINode
	p.poolRWMutex.Unlock()
}

func (p *INodePool) AllocRawINode() INodeUintptr {
	var uINode INodeUintptr
	uINode = INodeUintptr(p.inodeObjectPool.AllocRawObject())
	util.InitUUID64(&uINode.Ptr().ID)
	return uINode
}

func (p *INodePool) ReleaseRawINode(uINode INodeUintptr) {
	p.inodeObjectPool.ReleaseRawObject(uintptr(uINode))
}
