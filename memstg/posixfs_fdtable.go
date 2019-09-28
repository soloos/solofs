package memstg

import "soloos/common/solofsapitypes"

func (p *PosixFs) FdTableAllocFd(fsINodeID solofsapitypes.FsINodeID) solofsapitypes.FsINodeFileHandlerID {
	return p.FdTable.AllocFd(fsINodeID)
}

func (p *PosixFs) FdTableGetFd(fdID solofsapitypes.FsINodeFileHandlerID) solofsapitypes.FsINodeFileHandler {
	return p.FdTable.GetFd(fdID)
}

func (p *PosixFs) FdTableFdAddAppendPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddAppendPosition(fdID, delta)
}

func (p *PosixFs) FdTableFdAddReadPosition(fdID solofsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddReadPosition(fdID, delta)
}
