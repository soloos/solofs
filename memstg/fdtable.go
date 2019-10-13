package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/util"
	"soloos/solodb/offheap"
)

type FdTable struct {
	fdIDsPool  offheap.NoGCUintptrPool
	FdsRWMutex util.RWMutex
	Fds        []solofstypes.FsINodeFileHandler
}

func (p *FdTable) Init() error {
	p.fdIDsPool.New = func() uintptr {
		var (
			fdID uintptr
			fd   solofstypes.FsINodeFileHandler
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

func (p *FdTable) AllocFd(fsINodeIno solofstypes.FsINodeIno) solofstypes.FsINodeFileHandlerID {
	var fdID = solofstypes.FsINodeFileHandlerID(p.fdIDsPool.Get())
	p.FdsRWMutex.RLock()
	p.Fds[fdID].FsINodeIno = fsINodeIno
	p.FdsRWMutex.RUnlock()
	return fdID

}

func (p *FdTable) FdAddAppendPosition(fdID solofstypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].AppendPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) FdAddReadPosition(fdID solofstypes.FsINodeFileHandlerID, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].ReadPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) GetFd(fdID solofstypes.FsINodeFileHandlerID) (ret solofstypes.FsINodeFileHandler) {
	p.FdsRWMutex.RLock()
	ret = p.Fds[int(fdID)]
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) ReleaseFd(fdID solofstypes.FsINodeFileHandlerID) {
	p.Fds[int(fdID)].Reset()
	p.fdIDsPool.Put(uintptr(fdID))
}
