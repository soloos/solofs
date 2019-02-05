package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) StatLimits() (uint64, uint64) {
	var (
		capacity uint64 = 1024 * 1024 * 1024 * 1024 * 1024 * 100
		files    uint64 = 1024 * 1024 * 1024 * 100
	)
	return capacity, files
}

func (p *DirTreeStg) BlkSize() uint32 {
	// TODO return real result
	var (
		blksize uint32 = 1024 * 4
	)
	return blksize
}

func (p *DirTreeStg) StatFs(input *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
	// TODO return real result
	var (
		usedSize   uint64 = 1024 * 1024 * 100
		filesCount uint64 = 1000
	)

	capacity, files := p.StatLimits()
	blksize := p.BlkSize()

	out.Blocks = capacity / uint64(blksize)
	out.Bfree = (capacity - usedSize) / uint64(blksize)
	out.Bavail = (capacity - usedSize) / uint64(blksize)
	out.Files = files
	out.Ffree = out.Files - filesCount
	out.Bsize = blksize
	out.NameLen = types.FS_MAX_NAME_LENGTH
	out.Frsize = blksize

	return fuse.OK
}