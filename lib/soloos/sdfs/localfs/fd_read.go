package localfs

import (
	"soloos/sdfs/types"
)

func (p *Fd) PReadMemBlock(uMemBlock types.MemBlockUintptr,
	memBlockReadOffset int,
	memBlockReadEnd int,
	netINodeOffset int64,
) error {
	return p.ReadAt((*uMemBlock.Ptr().BytesSlice())[memBlockReadOffset:memBlockReadEnd], netINodeOffset)
}

func (p *Fd) ReadAt(data []byte, netINodeOffset int64) error {
	var (
		off int
		n   int
		err error
	)
	for off = 0; off < len(data); off += n {
		n, err = p.file.ReadAt(data, netINodeOffset+int64(off))
		if err != nil {
			return err
		}
	}
	return nil
}
