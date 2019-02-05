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

func (p *NetBlockPool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.netBlockObjectPool,
		int(NetBlockStructSize), -1,
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
		isLoaded        bool
	)

	EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)

	for {
		p.poolRWMutex.RLock()
		uNetBlock, exists = p.pool[netINodeBlockID]
		p.poolRWMutex.RUnlock()
		if exists {
			isLoaded = true
			goto FETCH_NETBLOCK_DONE
		}

		p.poolRWMutex.Lock()
		uNetBlock, exists = p.pool[netINodeBlockID]
		if exists {
			isLoaded = true
		} else {
			uNetBlock = NetBlockUintptr(p.netBlockObjectPool.AllocRawObject())
			uNetBlock.Ptr().NetINodeID = uNetINode.Ptr().ID
			uNetBlock.Ptr().IndexInNetINode = netBlockIndex
			uNetBlock.Ptr().SharedPointer.CompleteInit()
			isLoaded = false
			p.pool[netINodeBlockID] = uNetBlock
		}
		p.poolRWMutex.Unlock()

	FETCH_NETBLOCK_DONE:
		uNetBlock.Ptr().SharedPointer.ReadAcquire()

		if uNetBlock.Ptr().SharedPointer.IsInited() == false {
			uNetBlock.Ptr().SharedPointer.ReadRelease()
		} else {
			break
		}
	}

	return uNetBlock, isLoaded
}

func (p *NetBlockPool) ReleaseNetBlock(uNetBlock NetBlockUintptr) {
	pNetBlock := uNetBlock.Ptr()
	pNetBlock.SharedPointer.ReadRelease()
	if pNetBlock.SharedPointer.IsShouldRelease() &&
		pNetBlock.SharedPointer.Accessor == 0 {

		var netINodeBlockID NetINodeBlockID
		EncodeNetINodeBlockID(&netINodeBlockID, pNetBlock.NetINodeID, pNetBlock.IndexInNetINode)

		pNetBlock.SharedPointer.WriteAcquire()
		pNetBlock.Reset()
		p.poolRWMutex.Lock()
		delete(p.pool, netINodeBlockID)
		p.poolRWMutex.Unlock()
		pNetBlock.SharedPointer.WriteRelease()

		p.netBlockObjectPool.ReleaseRawObject(uintptr(uNetBlock))
	}
}
