package types

import (
	"soloos/common/util/offheap"
	"strconv"
	"sync"
)

type FsINodePool struct {
	offheapDriver     *offheap.OffheapDriver
	fsINodeObjectPool offheap.RawObjectPool

	fsINodesRWMutex sync.RWMutex
	fsINodesByID    map[FsINodeID]FsINodeUintptr
	fsINodesByPath  map[string]FsINodeUintptr
}

func (p *FsINodePool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.fsINodeObjectPool,
		int(FsINodeStructSize), -1,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.fsINodesByID = make(map[FsINodeID]FsINodeUintptr)
	p.fsINodesByPath = make(map[string]FsINodeUintptr)

	return nil
}

func (p *FsINodePool) MakeFsINodeKey(parentID FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + fsINodeName
}

func (p *FsINodePool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *FsINodePool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// return true if FsINode stored in fsINodesByPath before
//    	  false if FsINode is alloc
func (p *FsINodePool) MustGetFsINodeByPathWithReadAcquire(parentID FsINodeID, fsINodeName string) (FsINodeUintptr, bool) {
	var (
		fsINodeKey = p.MakeFsINodeKey(parentID, fsINodeName)
		uFsINode   FsINodeUintptr
		exists     bool
		isLoaded   bool
	)

	for {
		p.fsINodesRWMutex.RLock()
		uFsINode, exists = p.fsINodesByPath[fsINodeKey]
		p.fsINodesRWMutex.RUnlock()
		if exists {
			isLoaded = true
			goto FETCH_NETINODE_DONE
		}

		p.fsINodesRWMutex.Lock()
		uFsINode, exists = p.fsINodesByPath[fsINodeKey]
		if exists {
			isLoaded = true
		} else {
			uFsINode = FsINodeUintptr(p.fsINodeObjectPool.AllocRawObject())
			uFsINode.Ptr().ParentID = parentID
			uFsINode.Ptr().Name = fsINodeName
			uFsINode.Ptr().SharedPointer.CompleteInit()
			isLoaded = false
			p.fsINodesByPath[fsINodeKey] = uFsINode
		}
		p.fsINodesRWMutex.Unlock()

	FETCH_NETINODE_DONE:
		uFsINode.Ptr().SharedPointer.ReadAcquire()

		if uFsINode.Ptr().SharedPointer.IsInited() == false {
			uFsINode.Ptr().SharedPointer.ReadRelease()
		} else {
			break
		}
	}

	return uFsINode, isLoaded
}

// return true if FsINode stored in fsINodesByID before
//    	  false if FsINode is alloc
func (p *FsINodePool) MustGetFsINodeByIDWithReadAcquire(fsINodeID FsINodeID) (FsINodeUintptr, bool) {
	var (
		uFsINode FsINodeUintptr
		exists   bool
		isLoaded bool
	)

	for {
		p.fsINodesRWMutex.RLock()
		uFsINode, exists = p.fsINodesByID[fsINodeID]
		p.fsINodesRWMutex.RUnlock()
		if exists {
			isLoaded = true
			goto FETCH_NETINODE_DONE
		}

		p.fsINodesRWMutex.Lock()
		uFsINode, exists = p.fsINodesByID[fsINodeID]
		if exists {
			isLoaded = true
		} else {
			uFsINode = FsINodeUintptr(p.fsINodeObjectPool.AllocRawObject())
			uFsINode.Ptr().Ino = fsINodeID
			uFsINode.Ptr().SharedPointer.CompleteInit()
			isLoaded = false
			p.fsINodesByID[fsINodeID] = uFsINode
		}
		p.fsINodesRWMutex.Unlock()

	FETCH_NETINODE_DONE:
		uFsINode.Ptr().SharedPointer.ReadAcquire()

		if uFsINode.Ptr().SharedPointer.IsInited() == false {
			uFsINode.Ptr().SharedPointer.ReadRelease()
		} else {
			break
		}
	}

	return uFsINode, isLoaded
}

func (p *FsINodePool) ReleaseFsINodeWithReadRelease(uFsINode FsINodeUintptr) {
	if uFsINode == 0 {
		return
	}

	pFsINode := uFsINode.Ptr()
	pFsINode.SharedPointer.ReadRelease()
	if pFsINode.SharedPointer.IsShouldRelease() &&
		pFsINode.SharedPointer.Accessor == 0 {

		pFsINode.SharedPointer.WriteAcquire()
		pFsINode.Reset()
		p.fsINodesRWMutex.Lock()
		delete(p.fsINodesByID, pFsINode.Ino)
		delete(p.fsINodesByPath, p.MakeFsINodeKey(pFsINode.ParentID, pFsINode.Name))
		p.fsINodesRWMutex.Unlock()
		pFsINode.SharedPointer.WriteRelease()

		p.fsINodeObjectPool.ReleaseRawObject(uintptr(uFsINode))
	}
}
