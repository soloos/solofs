package types

import (
	"soloos/util/offheap"
	"sync"
)

type NetBlockPool struct {
	offheapDriver      *offheap.OffheapDriver
	netBlockObjectPool offheap.RawObjectPool
	poolRWMutex        sync.RWMutex
	pool               map[NetINodeBlockID]NetBlockUintptr
}

func (p *NetBlockPool) Init(rawChunksLimit int32,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.netBlockObjectPool,
		int(NetBlockStructSize), rawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[NetINodeBlockID]NetBlockUintptr)

	return nil
}

func (p *NetBlockPool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// return true if NetBlock stored in pool before
//    	  false if NetBlock is alloc
func (p *NetBlockPool) MustGetNetBlock(uNetINode NetINodeUintptr,
	netBlockIndex int) (NetBlockUintptr, bool) {
	var (
		netINodeBlockID NetINodeBlockID
		uNetBlock       NetBlockUintptr
		exists          bool
	)
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.poolRWMutex.RLock()
	uNetBlock, exists = p.pool[netINodeBlockID]
	p.poolRWMutex.RUnlock()
	if exists {
		return uNetBlock, true
	}

	var isLoaded = true

	p.poolRWMutex.Lock()
	uNetBlock, exists = p.pool[netINodeBlockID]
	if exists == false {
		uNetBlock = p.AllocRawNetBlock()
		uNetBlock.Ptr().NetINodeID = uNetINode.Ptr().ID
		uNetBlock.Ptr().IndexInNetINode = netBlockIndex
		isLoaded = false
		p.pool[netINodeBlockID] = uNetBlock
	}
	p.poolRWMutex.Unlock()

	return uNetBlock, isLoaded
}

// return true if set RawNetBlock success
//        false if there is RawNetBlock exists, and return the old NetBlock and release the new one
func (p *NetBlockPool) SetRawNetBlock(uNetINode NetINodeUintptr,
	netBlockIndex int,
	uNetBlock NetBlockUintptr) (NetBlockUintptr, bool) {
	var (
		netINodeBlockID NetINodeBlockID
		uNetBlockFinal  NetBlockUintptr
		exists          bool
	)
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.poolRWMutex.Lock()
	uNetBlockFinal, exists = p.pool[netINodeBlockID]
	if exists {
		p.ReleaseRawNetBlock(uNetBlock)
	} else {
		p.pool[netINodeBlockID] = uNetBlock
		uNetBlockFinal = uNetBlock
	}
	p.poolRWMutex.Unlock()
	return uNetBlockFinal, exists
}

func (p *NetBlockPool) AllocRawNetBlock() NetBlockUintptr {
	return NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
}

func (p *NetBlockPool) ReleaseRawNetBlock(uNetBlock NetBlockUintptr) {
	uNetBlock.Ptr().Reset()
	p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
}
