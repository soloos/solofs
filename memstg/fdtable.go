package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/util"
	"soloos/solodb/offheap"
)

type FdTable struct {
	fdIDsPool  offheap.NoGCUintptrPool
	FdsRWMutex util.RWMutex
	Fds        []solofsapitypes.FsINodeFileHandler
}

func (p *FdTable) Init() error {
	p.fdIDsPool.New = func() uintptr {
		var (
			fdID uintptr
			fd   solofsapitypes.FsINodeFileHandler
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

func (p *FdTable) AllocFd(fsINodeID solofsapitypes.FsINodeID) solofsapitypes.FsINodeFileHandlerID {
	var fdID = solofsapitypes.FsINodeFileHandlerID(p.fdIDsPool.Get())
	p.FdsRWMutex.RLock()
	p.Fds[fdID].FsINodeID = fsINodeID
	p.FdsRWMutex.RUnlock()
	return fdID

}

func (p *FdTable) FdAddAppendPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].AppendPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) FdAddReadPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].ReadPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) GetFd(fdID solofsapitypes.FsINodeFileHandlerID) (ret solofsapitypes.FsINodeFileHandler) {
	p.FdsRWMutex.RLock()
	ret = p.Fds[int(fdID)]
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) ReleaseFd(fdID solofsapitypes.FsINodeFileHandlerID) {
	p.Fds[int(fdID)].Reset()
	p.fdIDsPool.Put(uintptr(fdID))
}
