package memstg

import (
	"soloos/sdfs/types"
	"sync"
)

type FdTable struct {
	fdIDsPool  sync.NoGCUintptrPool
	FdsRWMutex sync.RWMutex
	Fds        []types.FsINodeFileHandler
}

func (p *FdTable) Init() error {
	p.fdIDsPool.New = func() uintptr {
		var (
			fdID uintptr
			fd   types.FsINodeFileHandler
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

func (p *FdTable) AllocFd(fsINodeID types.FsINodeID) uint64 {
	var fdID = int(p.fdIDsPool.Get())
	p.FdsRWMutex.RLock()
	p.Fds[fdID].FsINodeID = fsINodeID
	p.FdsRWMutex.RUnlock()
	return uint64(fdID)

}

func (p *FdTable) FdAddAppendPosition(fdID uint64, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].AppendPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) FdAddReadPosition(fdID uint64, delta uint64) {
	p.FdsRWMutex.RLock()
	p.Fds[int(fdID)].ReadPosition += delta
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) GetFd(fdID uint64) (ret types.FsINodeFileHandler) {
	p.FdsRWMutex.RLock()
	ret = p.Fds[int(fdID)]
	p.FdsRWMutex.RUnlock()
	return
}

func (p *FdTable) ReleaseFd(fdID uint64) {
	p.Fds[int(fdID)].Reset()
	p.fdIDsPool.Put(uintptr(fdID))
}
