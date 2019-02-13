package types

import (
	"soloos/common/util/offheap"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	INodeRWMutexStructSize = unsafe.Sizeof(INodeRWMutex{})
)

type INodeRWMutexID = FsINodeID
type INodeRWMutexUintptr uintptr

func (u INodeRWMutexUintptr) Ptr() *INodeRWMutex { return (*INodeRWMutex)(unsafe.Pointer(u)) }

type INodeRWMutexMeta struct {
	Start uint64
	End   uint64
	Typ   uint32
	Pid   uint32
}

type INodeRWMutex struct {
	SharedPointer offheap.SharedPointer `db:"-"`
	lockSigGuard  int32
	LockSig       sync.RWMutex

	ID INodeRWMutexID
	INodeRWMutexMeta
}

func (p *INodeRWMutex) Reset() {
	p.INodeRWMutexMeta.Start = 0
	p.INodeRWMutexMeta.End = 0
	p.INodeRWMutexMeta.Typ = 0
	p.INodeRWMutexMeta.Pid = 0
}

func (p *INodeRWMutex) RLock() {
	atomic.AddInt32(&p.lockSigGuard, 1)
	p.LockSig.RLock()
}

func (p *INodeRWMutex) RUnlock() {
	p.LockSig.RUnlock()
	atomic.AddInt32(&p.lockSigGuard, -1)
}

func (p *INodeRWMutex) Lock() {
	atomic.AddInt32(&p.lockSigGuard, 1)
	p.LockSig.Lock()
}

func (p *INodeRWMutex) Unlock() {
	p.LockSig.Unlock()
	atomic.AddInt32(&p.lockSigGuard, -1)
}

func (p *INodeRWMutex) TryRLock() {
	if atomic.AddInt32(&p.lockSigGuard, 1) == 1 {
		p.LockSig.RLock()
	} else {
		atomic.AddInt32(&p.lockSigGuard, -1)
	}
}

func (p *INodeRWMutex) TryLock() {
	if atomic.AddInt32(&p.lockSigGuard, 1) == 1 {
		p.LockSig.Lock()
	} else {
		atomic.AddInt32(&p.lockSigGuard, -1)
	}
}
