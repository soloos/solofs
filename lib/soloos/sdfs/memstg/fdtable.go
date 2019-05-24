package memstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/util"
	"soloos/sdbone/offheap"
)

type FdTable struct {
	fdIDsPool  offheap.NoGCUintptrPool
	FdsRWMutex util.RWMutex
	Fds        []sdfsapitypes.FsINodeFileHandler
}

func (p *FdTable) Init() error {
	p.fdIDsPool.New = func() uintptr {
		var (
			fdID uintptr
			fd   sdfsapitypes.FsINodeFileHandler
		)
		fd.Reset()
		p.FdsRWMutex.Lock()
		fdID = uintptr(len(p.Fds))
		p.Fds = append(p.Fds, fd)
		p.FdsRWMutex.Unlock()
		return fdID
	}
	return nil
}

func (p *FdTable) AllocFd(fsINodeID sdfsapitypes.FsINodeID) sdfsapitypes.FsINodeFileHandlerID {
	var fdID = sdfsapitypes.FsINodeFileHandlerID(p.fdIDsPool.Get())
	p.FdsRWMutex.RLock()
	p.Fds[fdID].FsINodeID = fsINodeID
	p.FdsRWMutex.RUnlock()
	return fdID

}

func (p *FdTable) FdAddAppendPosition(fdID sdfsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].AppendPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) FdAddReadPosition(fdID sdfsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].ReadPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) GetFd(fdID sdfsapitypes.FsINodeFileHandlerID) (ret sdfsapitypes.FsINodeFileHandler) {
	p.FdsRWMutex.RLock()
	ret = p.Fds[int(fdID)]
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) ReleaseFd(fdID sdfsapitypes.FsINodeFileHandlerID) {
	p.Fds[int(fdID)].Reset()
	p.fdIDsPool.Put(uintptr(fdID))
}
