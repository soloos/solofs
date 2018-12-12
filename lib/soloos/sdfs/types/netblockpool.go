package types

import (
	"soloos/util/offheap"
	"sync"
)

type NetBlockPool struct {
	offheapDriver      *offheap.OffheapDriver
	netBlockObjectPool offheap.RawObjectPool
	PoolRWMutex        sync.RWMutex
	Pool               map[NetINodeBlockID]NetBlockUintptr
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

	p.Pool = make(map[NetINodeBlockID]NetBlockUintptr)

	return nil
}

func (p *NetBlockPool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetNetBlock get or init a netBlockblock
func (p *NetBlockPool) MustGetNetBlock(uNetINode NetINodeUintptr, netBlockIndex int) (NetBlockUintptr, bool) {
	var (
		netINodeBlockID NetINodeBlockID
		uNetBlock       NetBlockUintptr
		exists          bool
	)

	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)

	p.PoolRWMutex.RLock()
	uNetBlock, exists = p.Pool[netINodeBlockID]
	p.PoolRWMutex.RUnlock()
	if exists {
		return uNetBlock, true
	}

	p.PoolRWMutex.Lock()
	uNetBlock, exists = p.Pool[netINodeBlockID]
	if exists {
		goto GET_DONE
	}

	uNetBlock = NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
	p.Pool[netINodeBlockID] = uNetBlock

GET_DONE:
	p.PoolRWMutex.Unlock()
	return uNetBlock, exists
}

func (p *NetBlockPool) ReleaseNetBlock(uNetINode NetINodeUintptr, netBlockIndex int, uNetBlock NetBlockUintptr) {
	var netINodeBlockID NetINodeBlockID
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.PoolRWMutex.Lock()
	delete(p.Pool, netINodeBlockID)
	p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
	p.PoolRWMutex.Unlock()
}

func (p *NetBlockPool) SetNetBlock(uNetINode NetINodeUintptr, netBlockIndex int, uNetBlock NetBlockUintptr) {
	var netINodeBlockID NetINodeBlockID
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.PoolRWMutex.Lock()
	p.Pool[netINodeBlockID] = uNetBlock
	p.PoolRWMutex.Unlock()
}

func (p *NetBlockPool) AllocRawNetBlock() NetBlockUintptr {
	return NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
}

func (p *NetBlockPool) ReleaseRawNetBlock(uNetBlock NetBlockUintptr) {
	p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
}
