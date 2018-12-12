package types

import (
	"soloos/util"
	"soloos/util/offheap"
	"sync"
)

type NetINodePool struct {
	offheapDriver      *offheap.OffheapDriver
	netINodeObjectPool offheap.RawObjectPool
	PoolRWMutex        sync.RWMutex
	Pool               map[NetINodeID]NetINodeUintptr
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

	p.Pool = make(map[NetINodeID]NetINodeUintptr)

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
		exists    bool
	)

	p.PoolRWMutex.RLock()
	uNetINode, exists = p.Pool[netINodeID]
	p.PoolRWMutex.RUnlock()
	if exists {
		return uNetINode, true
	}

	p.PoolRWMutex.Lock()
	uNetINode, exists = p.Pool[netINodeID]
	if exists {
		goto GET_DONE
	}

	uNetINode = NetINodeUintptr(p.netINodeObjectPool.AllocRawObject())
	uNetINode.Ptr().ID = netINodeID
	p.Pool[netINodeID] = uNetINode

GET_DONE:
	p.PoolRWMutex.Unlock()
	return uNetINode, exists
}

func (p *NetINodePool) ReleaseNetINode(uNetINode NetINodeUintptr) {
	var exists bool
	p.PoolRWMutex.Lock()
	_, exists = p.Pool[uNetINode.Ptr().ID]
	if exists {
		delete(p.Pool, uNetINode.Ptr().ID)
		p.netINodeObjectPool.ReleaseRawObject(uintptr(uNetINode))
	}
	p.PoolRWMutex.Unlock()
}

func (p *NetINodePool) SaveRawNetINode(uNetINode NetINodeUintptr) {
	p.PoolRWMutex.Lock()
	p.Pool[uNetINode.Ptr().ID] = uNetINode
	p.PoolRWMutex.Unlock()
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
