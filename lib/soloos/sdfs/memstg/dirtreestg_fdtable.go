package memstg

import "soloos/sdfs/types"

func (p *DirTreeStg) FdTableAllocFd(fsINodeID types.FsINodeID) uint64 {
	return p.FdTable.AllocFd(fsINodeID)
}

func (p *DirTreeStg) FdTableGetFd(fdID uint64) types.FsINodeFileHandler {
	return p.FdTable.GetFd(fdID)
}

func (p *DirTreeStg) FdTableFdAddAppendPosition(fdID uint64, delta uint64) {
	p.FdTable.FdAddAppendPosition(fdID, delta)
}

func (p *DirTreeStg) FdTableFdAddReadPosition(fdID uint64, delta uint64) {
	p.FdTable.FdAddReadPosition(fdID, delta)
}
