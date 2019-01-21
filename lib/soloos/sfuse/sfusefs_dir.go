package sfuse

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *SFuseFs) Mkdir(input *fuse.MkdirIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.Client.MemDirTreeStg.Mkdir(nil, input.NodeId, input.Mode, name, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFuseEntryOutByFsINode(out, &fsINode)
	return fuse.OK
}

func (p *SFuseFs) Rmdir(header *fuse.InHeader, name string) (code fuse.Status) {
	var err error
	err = p.Client.MemDirTreeStg.Rmdir(header.NodeId)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}
	return fuse.OK
}

// Directory handling
func (p *SFuseFs) OpenDir(input *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByID(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	out.Fh = p.Client.MemDirTreeStg.FdTable.AllocFd(fsINode.Ino)
	out.OpenFlags = input.Flags
	return fuse.OK
}

func (p *SFuseFs) ReadDir(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		isAddDirEntrySuccess bool
		err                  error
	)
	err = p.Client.MemDirTreeStg.ListFsINodeByIno(input.NodeId, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount) - input.Offset, input.Offset
		},
		func(fsINode types.FsINode) bool {
			isAddDirEntrySuccess, _ = out.AddDirEntry(fuse.DirEntry{
				Mode: fsINode.Mode,
				Name: fsINode.Name,
				Ino:  fsINode.Ino,
			})
			return isAddDirEntrySuccess
		},
	)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *SFuseFs) ReadDirPlus(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		entryOut *fuse.EntryOut
		off      uint64
		err      error
	)
	err = p.Client.MemDirTreeStg.ListFsINodeByIno(input.NodeId, true,
		func(resultCount int) (uint64, uint64) {
			var fetchRowsLimit uint64
			if uint64(resultCount) > input.Offset {
				fetchRowsLimit = uint64(resultCount) - input.Offset
				if fetchRowsLimit > 1024 {
					fetchRowsLimit = 1024
				}
			} else {
				fetchRowsLimit = 0
			}
			return fetchRowsLimit, input.Offset
		},
		func(fsINode types.FsINode) bool {
			entryOut, off = out.AddDirLookupEntry(fuse.DirEntry{
				Mode: fsINode.Mode,
				Name: fsINode.Name,
				Ino:  fsINode.Ino,
			})
			if entryOut == nil {
				return false
			}
			p.setFuseEntryOutByFsINode(entryOut, &fsINode)
			return true
		},
	)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}
	return fuse.OK
}

func (p *SFuseFs) ReleaseDir(input *fuse.ReleaseIn) {
	// TODO make sure releaable
	p.Client.MemDirTreeStg.FdTable.ReleaseFd(input.Fh)
}

func (p *SFuseFs) FsyncDir(input *fuse.FsyncIn) (code fuse.Status) {
	return fuse.OK
}
