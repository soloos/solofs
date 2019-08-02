package memstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/sdbone/offheap"
	"soloos/sdfs/sdfstypes"
	"sync/atomic"
)

// TODO release sdfsapitypes.INodeRWMutex
// support distributed system

// GetLk returns existing lock information for file
func (p *FsINodeDriver) GetLk(fsINodeID sdfsapitypes.FsINodeID, iNodeRWMutexMeta *sdfsapitypes.INodeRWMutexMeta) error {
	var (
		uObject       offheap.HKVTableObjectUPtrWithUint64
		uINodeRWMutex sdfsapitypes.INodeRWMutexUintptr
	)

	uObject, _ = p.INodeRWMutexTable.MustGetObjectWithReadAcquire(fsINodeID)
	uINodeRWMutex = sdfsapitypes.INodeRWMutexUintptr(uObject)
	defer p.INodeRWMutexTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithUint64(uINodeRWMutex))
	*iNodeRWMutexMeta = uINodeRWMutex.Ptr().INodeRWMutexMeta
	return nil
}

func (p *FsINodeDriver) doSetLk(fsINodeID sdfsapitypes.FsINodeID, setFlag *sdfsapitypes.INodeRWMutexMeta, isShouldBlock bool) error {
	var (
		uObject       offheap.HKVTableObjectUPtrWithUint64
		uINodeRWMutex sdfsapitypes.INodeRWMutexUintptr
		pINodeRWMutex *sdfsapitypes.INodeRWMutex
		err           error
	)

	uObject, _ = p.INodeRWMutexTable.MustGetObjectWithReadAcquire(fsINodeID)
	uINodeRWMutex = sdfsapitypes.INodeRWMutexUintptr(uObject)
	defer p.INodeRWMutexTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithUint64(uINodeRWMutex))

	pINodeRWMutex = uINodeRWMutex.Ptr()

	pINodeRWMutex.INodeRWMutexMeta.Start = setFlag.Start
	pINodeRWMutex.INodeRWMutexMeta.End = setFlag.End
	pINodeRWMutex.INodeRWMutexMeta.Pid = setFlag.Pid

	if setFlag.Typ == sdfstypes.FS_INODE_LOCK_SH {
		if isShouldBlock {
			pINodeRWMutex.RLock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = sdfstypes.FS_INODE_LOCK_SH
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(sdfstypes.FS_INODE_LOCK_SH)) ||
				pINodeRWMutex.INodeRWMutexMeta.Typ == sdfstypes.FS_INODE_LOCK_SH {
				go pINodeRWMutex.RLock()
				err = nil
			} else {
				err = sdfsapitypes.ErrRLockFailed
			}
		}

	} else if setFlag.Typ == sdfstypes.FS_INODE_LOCK_EX {
		if isShouldBlock {
			pINodeRWMutex.Lock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = sdfstypes.FS_INODE_LOCK_EX
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(sdfstypes.FS_INODE_LOCK_EX)) {
				go pINodeRWMutex.LockSig.Lock()
				err = nil
			} else {
				err = sdfsapitypes.ErrLockFailed
			}
		}

	} else {
		err = sdfsapitypes.ErrInvalidArgs
	}

	return err
}

// SetLk Sets or clears the lock described by lk on file.
func (p *FsINodeDriver) SetLk(fsINodeID sdfsapitypes.FsINodeID, setFlag *sdfsapitypes.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, false)
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *FsINodeDriver) SetLkw(fsINodeID sdfsapitypes.FsINodeID, setFlag *sdfsapitypes.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, true)
}
