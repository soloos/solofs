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

func (p *NetINodePool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.netINodeObjectPool,
		int(NetINodeStructSize), -1,
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
	uNetINode := NetINodeUintptr(uRawChunk)
	uNetINode.Ptr().SharedPointer.IsInited = true
}

// return true if NetINode stored in pool before
//    	  false if NetINode is alloc
func (p *NetINodePool) MustGetNetINodeWithReadAcquire(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	var (
		uNetINode NetINodeUintptr
		exists    bool
		isLoaded  bool
	)

	for {
		p.poolRWMutex.RLock()
		uNetINode, exists = p.pool[netINodeID]
		p.poolRWMutex.RUnlock()
		if exists {
			isLoaded = true
			goto FETCH_NETINODE_DONE
		}

		p.poolRWMutex.Lock()
		uNetINode, exists = p.pool[netINodeID]
		if exists {
			isLoaded = true
		} else {
			uNetINode = NetINodeUintptr(p.netINodeObjectPool.AllocRawObject())
			uNetINode.Ptr().ID = netINodeID
			isLoaded = false
			p.pool[netINodeID] = uNetINode
		}
		p.poolRWMutex.Unlock()

	FETCH_NETINODE_DONE:
		uNetINode.Ptr().SharedPointer.ReadAcquire()

		if uNetINode.Ptr().SharedPointer.IsInited == false {
			uNetINode.Ptr().SharedPointer.ReadRelease()
		} else {
			break
		}
	}

	return uNetINode, isLoaded
}

func (p *NetINodePool) ReleaseNetINodeWithReadRelease(uNetINode NetINodeUintptr) {
	if uNetINode == 0 {
		return
	}

	pNetINode := uNetINode.Ptr()
	pNetINode.SharedPointer.ReadRelease()
	if pNetINode.SharedPointer.IsShouldRelease &&
		pNetINode.SharedPointer.Accessor == 0 {

		pNetINode.SharedPointer.WriteAcquire()
		pNetINode.Reset()
		p.poolRWMutex.Lock()
		delete(p.pool, pNetINode.ID)
		p.poolRWMutex.Unlock()
		pNetINode.SharedPointer.WriteRelease()

		p.netINodeObjectPool.ReleaseRawObject(uintptr(uNetINode))
	}
}
