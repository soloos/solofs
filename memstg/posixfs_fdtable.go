package memstg

import "soloos/common/solofstypes"

func (p *PosixFs) FdTableAllocFd(fsINodeID solofstypes.FsINodeID) solofstypes.FsINodeFileHandlerID {
	return p.FdTable.AllocFd(fsINodeID)
}

func (p *PosixFs) FdTableGetFd(fdID solofstypes.FsINodeFileHandlerID) solofstypes.FsINodeFileHandler {
	return p.FdTable.GetFd(fdID)
}

func (p *PosixFs) FdTableFdAddAppendPosition(fdID solofstypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddAppendPosition(fdID, delta)
}

func (p *PosixFs) FdTableFdAddReadPosition(fdID solofstypes.FsINodeFileHandlerID, delta uint64) {
	p.FdTable.FdAddReadPosition(fdID, delta)
}
