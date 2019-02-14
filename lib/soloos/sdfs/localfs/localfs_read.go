package localfs

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

func (p *LocalFs) PReadMemBlockWithDisk(uNetINode types.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int32,
	uMemBlock types.MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error) {
	var (
		fd                 *Fd
		memBlockReadOffset int
		readedLen          int
		err                error
	)

	fd, err = p.fdDriver.Open(uNetINode, uNetBlock)
	if err != nil {
		goto PREAD_DONE
	}

	memBlockReadOffset = int(offset - uint64(memBlockIndex)*uint64(uMemBlock.Ptr().Bytes.Cap))
	readedLen, err = fd.PReadMemBlock(uMemBlock,
		memBlockReadOffset,
		memBlockReadOffset+length,
		offset)
	if err != nil {
		goto PREAD_DONE
	}

PREAD_DONE:
	// TODO catch close file error
	p.fdDriver.Close(fd)

	return readedLen, nil
}
