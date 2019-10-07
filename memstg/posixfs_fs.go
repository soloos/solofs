package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofstypes"
)

func (p *PosixFs) StatLimits() (uint64, uint64) {
	var (
		capacity uint64 = 1024 * 1024 * 1024 * 1024 * 1024 * 100
		files    uint64 = 1024 * 1024 * 1024 * 100
	)
	return capacity, files
}

func (p *PosixFs) BlkSize() uint32 {
	// TODO return real result
	var (
		blksize uint32 = 1024 * 4
	)
	return blksize
}

func (p *PosixFs) StatFs(input *fsapitypes.InHeader, out *fsapitypes.StatfsOut) (code fsapitypes.Status) {
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
	out.NameLen = solofstypes.FS_MAX_NAME_LENGTH
	out.Frsize = blksize

	return fsapitypes.OK
}
