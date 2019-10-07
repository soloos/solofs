package memstg

import (
	"soloos/common/solofstypes"
	"soloos/solodb/offheap"
	"sync/atomic"
)

// TODO release solofstypes.INodeRWMutex
// support distributed system

// GetLk returns existing lock information for file
func (p *FsMutexDriver) GetLk(fsINodeID solofstypes.FsINodeID, iNodeRWMutexMeta *solofstypes.INodeRWMutexMeta) error {
	var (
		uObject       offheap.LKVTableObjectUPtrWithUint64
		uINodeRWMutex solofstypes.INodeRWMutexUintptr
	)

	uObject, _ = p.INodeRWMutexTable.MustGetObject(fsINodeID)
	uINodeRWMutex = solofstypes.INodeRWMutexUintptr(uObject)
	defer p.INodeRWMutexTable.ReleaseObject(offheap.LKVTableObjectUPtrWithUint64(uINodeRWMutex))
	*iNodeRWMutexMeta = uINodeRWMutex.Ptr().INodeRWMutexMeta
	return nil
}

func (p *FsMutexDriver) doSetLk(fsINodeID solofstypes.FsINodeID, setFlag *solofstypes.INodeRWMutexMeta, isShouldBlock bool) error {
	var (
		uObject       offheap.LKVTableObjectUPtrWithUint64
		uINodeRWMutex solofstypes.INodeRWMutexUintptr
		pINodeRWMutex *solofstypes.INodeRWMutex
		err           error
	)

	uObject, _ = p.INodeRWMutexTable.MustGetObject(fsINodeID)
	uINodeRWMutex = solofstypes.INodeRWMutexUintptr(uObject)
	defer p.INodeRWMutexTable.ReleaseObject(offheap.LKVTableObjectUPtrWithUint64(uINodeRWMutex))

	pINodeRWMutex = uINodeRWMutex.Ptr()

	pINodeRWMutex.INodeRWMutexMeta.Start = setFlag.Start
	pINodeRWMutex.INodeRWMutexMeta.End = setFlag.End
	pINodeRWMutex.INodeRWMutexMeta.Pid = setFlag.Pid

	if setFlag.Typ == solofstypes.FS_INODE_LOCK_SH {
		if isShouldBlock {
			pINodeRWMutex.RLock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = solofstypes.FS_INODE_LOCK_SH
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(solofstypes.FS_INODE_LOCK_SH)) ||
				pINodeRWMutex.INodeRWMutexMeta.Typ == solofstypes.FS_INODE_LOCK_SH {
				go pINodeRWMutex.RLock()
				err = nil
			} else {
				err = solofstypes.ErrRLockFailed
			}
		}

	} else if setFlag.Typ == solofstypes.FS_INODE_LOCK_EX {
		if isShouldBlock {
			pINodeRWMutex.Lock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = solofstypes.FS_INODE_LOCK_EX
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(solofstypes.FS_INODE_LOCK_EX)) {
				go pINodeRWMutex.LockSig.Lock()
				err = nil
			} else {
				err = solofstypes.ErrLockFailed
			}
		}

	} else {
		err = solofstypes.ErrInvalidArgs
	}

	return err
}

// SetLk Sets or clears the lock described by lk on file.
func (p *FsMutexDriver) SetLk(fsINodeID solofstypes.FsINodeID, setFlag *solofstypes.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, false)
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *FsMutexDriver) SetLkw(fsINodeID solofstypes.FsINodeID, setFlag *solofstypes.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, true)
}
