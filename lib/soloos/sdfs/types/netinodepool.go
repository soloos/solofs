package types

import (
	"soloos/util/offheap"
	"sync"
)

type NetINodePool struct {
	offheapDriver      *offheap.OffheapDriver
	netINodeObjectPool offheap.RawObjectPool
	poolRWMutex        sync.RWMutex
	pool               map[NetINodeID]NetINodeUintptr
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

// return true if NetINode stored in pool before
//    	  false if NetINode is alloc
func (p *NetINodePool) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	var (
		uNetINode NetINodeUintptr
		exists    bool
	)
	p.poolRWMutex.RLock()
	uNetINode, exists = p.pool[netINodeID]
	p.poolRWMutex.RUnlock()
	if exists {
		return uNetINode, true
	}

	var isLoaded = true

	p.poolRWMutex.Lock()
	uNetINode, exists = p.pool[netINodeID]
	if exists == false {
		uNetINode = p.AllocRawNetINode()
		uNetINode.Ptr().ID = netINodeID
		isLoaded = false
		p.pool[netINodeID] = uNetINode
	}
	p.poolRWMutex.Unlock()

	return uNetINode, isLoaded
}

// return true if set RawNetINode success
//        false if there is RawNetINode exists, and return the old NetINode and release the new one
func (p *NetINodePool) SetRawNetINode(uNetINode NetINodeUintptr) (NetINodeUintptr, bool) {
	var (
		uNetINodeFinal NetINodeUintptr
		exists         bool
	)
	p.poolRWMutex.Lock()
	uNetINodeFinal, exists = p.pool[uNetINode.Ptr().ID]
	if exists {
		p.ReleaseRawNetINode(uNetINode)
	} else {
		p.pool[uNetINode.Ptr().ID] = uNetINode
		uNetINodeFinal = uNetINode
	}
	p.poolRWMutex.Unlock()
	return uNetINodeFinal, exists
}

func (p *NetINodePool) AllocRawNetINode() NetINodeUintptr {
	var uNetINode NetINodeUintptr
	uNetINode = NetINodeUintptr(p.netINodeObjectPool.AllocRawObject())
	return uNetINode
}

func (p *NetINodePool) ReleaseRawNetINode(uNetINode NetINodeUintptr) {
	uNetINode.Ptr().Reset()
	p.netINodeObjectPool.ReleaseRawObject(uintptr(uNetINode))
}
