package types

import (
	"soloos/util"
	"soloos/util/offheap"
	"sync"
)

type NetINodePool struct {
	offheapDriver   *offheap.OffheapDriver
	netINodeObjectPool offheap.RawObjectPool
	poolRWMutex     sync.RWMutex
	pool            map[NetINodeID]NetINodeUintptr
}

func (p *NetINodePool) Init(rawChunksLimit int32,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.netINodeObjectPool,
		int(NetINodeStructSize), rawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[NetINodeID]NetINodeUintptr)

	return nil
}

func (p *NetINodePool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetINodePool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetNetINode get or init a netINodeblock
func (p *NetINodePool) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	var (
		uNetINode NetINodeUintptr
		exists bool
	)

	p.poolRWMutex.RLock()
	uNetINode, exists = p.pool[netINodeID]
	p.poolRWMutex.RUnlock()
	if exists {
		return uNetINode, true
	}

	p.poolRWMutex.Lock()
	uNetINode, exists = p.pool[netINodeID]
	if exists {
		goto GET_DONE
	}

	uNetINode = NetINodeUintptr(p.netINodeObjectPool.AllocRawObject())
	uNetINode.Ptr().ID = netINodeID
	p.pool[netINodeID] = uNetINode

GET_DONE:
	p.poolRWMutex.Unlock()
	return uNetINode, exists
}

func (p *NetINodePool) ReleaseNetINode(uNetINode NetINodeUintptr) {
	var exists bool
	p.poolRWMutex.Lock()
	_, exists = p.pool[uNetINode.Ptr().ID]
	if exists {
		delete(p.pool, uNetINode.Ptr().ID)
		p.netINodeObjectPool.ReleaseRawObject(uintptr(uNetINode))
	}
	p.poolRWMutex.Unlock()
}

func (p *NetINodePool) SaveRawNetINode(uNetINode NetINodeUintptr) {
	p.poolRWMutex.Lock()
	p.pool[uNetINode.Ptr().ID] = uNetINode
	p.poolRWMutex.Unlock()
}

func (p *NetINodePool) AllocRawNetINode() NetINodeUintptr {
	var uNetINode NetINodeUintptr
	uNetINode = NetINodeUintptr(p.netINodeObjectPool.AllocRawObject())
	util.InitUUID64(&uNetINode.Ptr().ID)
	return uNetINode
}

func (p *NetINodePool) ReleaseRawNetINode(uNetINode NetINodeUintptr) {
	p.netINodeObjectPool.ReleaseRawObject(uintptr(uNetINode))
}
