package memstg

import (
	"io"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

type preadArg struct {
	conn       *snettypes.Connection
	dataLength int
	data       []byte
	offset     int64
}

func (p *NetINodeDriver) doPRead(uNetINode types.NetINodeUintptr,
	arg preadArg) (int, error) {
	var (
		uMemBlock          types.MemBlockUintptr
		uNetBlock          types.NetBlockUintptr
		memBlockIndex      int
		netBlockIndex      int
		memBlockStart      int64
		memBlockReadOffset int
		// memBlockReadEnd     int
		memBlockReadLength int
		offset             = arg.offset
		dataOffset         int
		readEnd            int64
		err                error
	)
	pNetINode := uNetINode.Ptr()

	if pNetINode.Size < arg.offset {
		return 0, io.EOF
	}

	if arg.offset+int64(arg.dataLength) > pNetINode.Size {
		arg.dataLength = int(pNetINode.Size - arg.offset)
	}

	readEnd = offset + int64(arg.dataLength)
	for ; offset < readEnd; offset, dataOffset = offset+int64(memBlockReadLength), dataOffset+memBlockReadLength {
		// prepare netBlock
		netBlockIndex = int(offset / int64(pNetINode.NetBlockCap))
		uNetBlock, err = p.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)

		// prepare memBlock
		memBlockIndex = int(offset / int64(pNetINode.MemBlockCap))
		memBlockStart = int64(memBlockIndex) * int64(pNetINode.MemBlockCap)
		memBlockReadOffset = int(offset - memBlockStart)
		if memBlockStart+int64(pNetINode.MemBlockCap) < readEnd {
			// not the last block
			memBlockReadLength = int(memBlockStart + int64(pNetINode.MemBlockCap) - offset)
		} else {
			// the last block
			memBlockReadLength = int(readEnd - offset)
		}
		// memBlockReadEnd = memBlockReadOffset + memBlockReadLength
		uMemBlock, _ = p.memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)

		// TODO maybe rebase is not needed
		if uMemBlock.Ptr().Contains(memBlockReadOffset, memBlockReadLength) == false {
			err = p.unsafeMemBlockRebaseNetBlock(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex)
			if err != nil {
				goto READ_DATA_ONE_RUN_DONE
			}
		}

		// read memblock
		if arg.conn == nil {
			uMemBlock.Ptr().PReadWithMem(arg.data[dataOffset:dataOffset+memBlockReadLength], memBlockReadOffset)
		} else {
			// TODO read to connection
			err = uMemBlock.Ptr().PReadWithConn(arg.conn, memBlockReadLength, memBlockReadOffset)
			if err != nil {
				goto READ_DATA_ONE_RUN_DONE
			}
		}

	READ_DATA_ONE_RUN_DONE:
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
		if err != nil {
			goto READ_DATA_DONE
		}
	}

READ_DATA_DONE:
	return arg.dataLength, err
}

func (p *NetINodeDriver) PReadWithConn(uNetINode types.NetINodeUintptr,
	conn *snettypes.Connection, dataLength int, offset int64) (int, error) {
	return p.doPRead(uNetINode, preadArg{
		conn:       conn,
		data:       nil,
		dataLength: dataLength,
		offset:     offset,
	})
}

func (p *NetINodeDriver) PReadWithMem(uNetINode types.NetINodeUintptr,
	data []byte, offset int64) (int, error) {
	return p.doPRead(uNetINode, preadArg{
		conn:       nil,
		data:       data,
		dataLength: len(data),
		offset:     offset,
	})
}
