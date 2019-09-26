package memstg

import "soloos/common/solofsapitypes"

func (p *PosixFS) FdTableAllocFd(fsINodeID solofsapitypes.FsINodeID) solofsapitypes.FsINodeFileHandlerID {
	return p.FdTable.AllocFd(fsINodeID)
}

func (p *PosixFS) FdTableGetFd(fdID solofsapitypes.FsINodeFileHandlerID) solofsapitypes.FsINodeFileHandler {
	return p.FdTable.GetFd(fdID)
}

func (p *PosixFS) FdTableFdAddAppendPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddAppendPosition(fdID, delta)
}

func (p *PosixFS) FdTableFdAddReadPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddReadPosition(fdID, delta)
}
