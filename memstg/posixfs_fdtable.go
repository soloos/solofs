package memstg

import "soloos/common/sdfsapitypes"

func (p *PosixFS) FdTableAllocFd(fsINodeID sdfsapitypes.FsINodeID) sdfsapitypes.FsINodeFileHandlerID {
	return p.FdTable.AllocFd(fsINodeID)
}

func (p *PosixFS) FdTableGetFd(fdID sdfsapitypes.FsINodeFileHandlerID) sdfsapitypes.FsINodeFileHandler {
	return p.FdTable.GetFd(fdID)
}

func (p *PosixFS) FdTableFdAddAppendPosition(fdID sdfsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddAppendPosition(fdID, delta)
}

func (p *PosixFS) FdTableFdAddReadPosition(fdID sdfsapitypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddReadPosition(fdID, delta)
}
