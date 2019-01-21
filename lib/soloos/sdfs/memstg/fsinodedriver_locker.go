package memstg

import (
	"soloos/sdfs/types"
	"sync/atomic"
)

// TODO release types.INodeRWMutex
// support distributed system

// GetLk returns existing lock information for file
func (p *FsINodeDriver) GetLk(fsINodeID types.FsINodeID, iNodeRWMutexMeta *types.INodeRWMutexMeta) error {
	var (
		uINodeRWMutex types.INodeRWMutexUintptr
	)

	uINodeRWMutex, _ = p.INodeRWMutexPool.MustGetINodeRWMutexWithReadAcquire(fsINodeID)
	defer p.INodeRWMutexPool.ReleaseINodeRWMutexWithReadRelease(uINodeRWMutex)
	*iNodeRWMutexMeta = uINodeRWMutex.Ptr().INodeRWMutexMeta
	return nil
}

func (p *FsINodeDriver) doSetLk(fsINodeID types.FsINodeID, setFlag *types.INodeRWMutexMeta, isShouldBlock bool) error {
	var (
		uINodeRWMutex types.INodeRWMutexUintptr
		pINodeRWMutex *types.INodeRWMutex
		err           error
	)

	uINodeRWMutex, _ = p.INodeRWMutexPool.MustGetINodeRWMutexWithReadAcquire(fsINodeID)
	defer p.INodeRWMutexPool.ReleaseINodeRWMutexWithReadRelease(uINodeRWMutex)

	pINodeRWMutex = uINodeRWMutex.Ptr()

	pINodeRWMutex.INodeRWMutexMeta.Start = setFlag.Start
	pINodeRWMutex.INodeRWMutexMeta.End = setFlag.End
	pINodeRWMutex.INodeRWMutexMeta.Pid = setFlag.Pid

	if setFlag.Typ == types.FS_INODE_LOCK_SH {
		if isShouldBlock {
			pINodeRWMutex.RLock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = types.FS_INODE_LOCK_SH
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(types.FS_INODE_LOCK_SH)) ||
				pINodeRWMutex.INodeRWMutexMeta.Typ == types.FS_INODE_LOCK_SH {
				go pINodeRWMutex.RLock()
				err = nil
			} else {
				err = types.ErrRLockFailed
			}
		}

	} else if setFlag.Typ == types.FS_INODE_LOCK_EX {
		if isShouldBlock {
			pINodeRWMutex.Lock()
			pINodeRWMutex.INodeRWMutexMeta.Typ = types.FS_INODE_LOCK_EX
			err = nil
		} else {
			if atomic.CompareAndSwapUint32(&pINodeRWMutex.INodeRWMutexMeta.Typ, 0, uint32(types.FS_INODE_LOCK_EX)) {
				go pINodeRWMutex.LockSig.Lock()
				err = nil
			} else {
				err = types.ErrLockFailed
			}
		}

	} else {
		err = types.ErrInvalidArgs
	}

	return err
}

// SetLk Sets or clears the lock described by lk on file.
func (p *FsINodeDriver) SetLk(fsINodeID types.FsINodeID, setFlag *types.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, false)
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *FsINodeDriver) SetLkw(fsINodeID types.FsINodeID, setFlag *types.INodeRWMutexMeta) error {
	return p.doSetLk(fsINodeID, setFlag, true)
}
