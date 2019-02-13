package types

import (
	"soloos/common/util/offheap"
	"sync"
)

type INodeRWMutexPool struct {
	offheapDriver          *offheap.OffheapDriver
	inodeRWMutexObjectPool offheap.RawObjectPool
	poolRWMutex            sync.RWMutex
	pool                   map[INodeRWMutexID]INodeRWMutexUintptr
}

func (p *INodeRWMutexPool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.inodeRWMutexObjectPool,
		int(INodeRWMutexStructSize), -1,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[INodeRWMutexID]INodeRWMutexUintptr)

	return nil
}

func (p *INodeRWMutexPool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *INodeRWMutexPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// return true if INodeRWMutex stored in pool before
//    	  false if INodeRWMutex is alloc
func (p *INodeRWMutexPool) MustGetINodeRWMutexWithReadAcquire(inodeRWMutexID INodeRWMutexID) (INodeRWMutexUintptr, bool) {
	var (
		uINodeRWMutex INodeRWMutexUintptr
		exists        bool
		isLoaded      bool
	)

	for {
		p.poolRWMutex.RLock()
		uINodeRWMutex, exists = p.pool[inodeRWMutexID]
		p.poolRWMutex.RUnlock()
		if exists {
			isLoaded = true
			goto FETCH_NETINODE_DONE
		}

		p.poolRWMutex.Lock()
		uINodeRWMutex, exists = p.pool[inodeRWMutexID]
		if exists {
			isLoaded = true
		} else {
			uINodeRWMutex = INodeRWMutexUintptr(p.inodeRWMutexObjectPool.AllocRawObject())
			uINodeRWMutex.Ptr().ID = inodeRWMutexID
			uINodeRWMutex.Ptr().SharedPointer.CompleteInit()
			isLoaded = false
			p.pool[inodeRWMutexID] = uINodeRWMutex
		}
		p.poolRWMutex.Unlock()

	FETCH_NETINODE_DONE:
		uINodeRWMutex.Ptr().SharedPointer.ReadAcquire()

		if uINodeRWMutex.Ptr().SharedPointer.IsInited() == false {
			uINodeRWMutex.Ptr().SharedPointer.ReadRelease()
		} else {
			break
		}
	}

	return uINodeRWMutex, isLoaded
}

func (p *INodeRWMutexPool) ReleaseINodeRWMutexWithReadRelease(uINodeRWMutex INodeRWMutexUintptr) {
	if uINodeRWMutex == 0 {
		return
	}

	pINodeRWMutex := uINodeRWMutex.Ptr()
	pINodeRWMutex.SharedPointer.ReadRelease()
	if pINodeRWMutex.SharedPointer.IsShouldRelease() &&
		pINodeRWMutex.SharedPointer.Accessor == 0 {

		pINodeRWMutex.SharedPointer.WriteAcquire()
		pINodeRWMutex.Reset()
		p.poolRWMutex.Lock()
		delete(p.pool, pINodeRWMutex.ID)
		p.poolRWMutex.Unlock()
		pINodeRWMutex.SharedPointer.WriteRelease()

		p.inodeRWMutexObjectPool.ReleaseRawObject(uintptr(uINodeRWMutex))
	}
}
