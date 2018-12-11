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

// MustGetNetBlock get or init a netBlockblock
func (p *NetBlockPool) MustGetNetBlock(uNetINode NetINodeUintptr, netBlockIndex int) (NetBlockUintptr, bool) {
	var (
		netINodeBlockID NetINodeBlockID
		uNetBlock    NetBlockUintptr
		exists       bool
	)

	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)

	p.poolRWMutex.RLock()
	uNetBlock, exists = p.pool[netINodeBlockID]
	p.poolRWMutex.RUnlock()
	if exists {
		return uNetBlock, true
	}

	p.poolRWMutex.Lock()
	uNetBlock, exists = p.pool[netINodeBlockID]
	if exists {
		goto GET_DONE
	}

	uNetBlock = NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
	p.pool[netINodeBlockID] = uNetBlock

GET_DONE:
	p.poolRWMutex.Unlock()
	return uNetBlock, exists
}

func (p *NetBlockPool) ReleaseNetBlock(uNetINode NetINodeUintptr, netBlockIndex int, uNetBlock NetBlockUintptr) {
	var netINodeBlockID NetINodeBlockID
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.poolRWMutex.Lock()
	delete(p.pool, netINodeBlockID)
	p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
	p.poolRWMutex.Unlock()
}

func (p *NetBlockPool) SetNetBlock(uNetINode NetINodeUintptr, netBlockIndex int, uNetBlock NetBlockUintptr) {
	var netINodeBlockID NetINodeBlockID
	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	p.poolRWMutex.Lock()
	p.pool[netINodeBlockID] = uNetBlock
	p.poolRWMutex.Unlock()
}

func (p *NetBlockPool) AllocRawNetBlock() NetBlockUintptr {
	return NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
}

func (p *NetBlockPool) ReleaseRawNetBlock(uNetBlock NetBlockUintptr) {
	p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
}
