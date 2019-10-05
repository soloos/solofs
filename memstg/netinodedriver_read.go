package memstg

import (
	"io"
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

type preadArg struct {
	netQuery   *snettypes.NetQuery
	dataLength int
	data       []byte
	offset     uint64
}

func (p *NetINodeDriver) doPRead(uNetINode solofsapitypes.NetINodeUintptr,
	arg preadArg) (int, error) {
	var (
		uMemBlock          solofsapitypes.MemBlockUintptr
		uNetBlock          solofsapitypes.NetBlockUintptr
		memBlockIndex      int32
		netBlockIndex      int32
		memBlockStart      uint64
		memBlockReadOffset int
		// memBlockReadEnd     int
		memBlockReadLength int
		offset             = arg.offset
		dataOffset         int
		readEnd            uint64
		err                error
	)
	pNetINode := uNetINode.Ptr()

	if pNetINode.Size < arg.offset {
		return 0, io.EOF
	}

	if arg.offset+uint64(arg.dataLength) > pNetINode.Size {
		arg.dataLength = int(pNetINode.Size - arg.offset)
	}

	readEnd = offset + uint64(arg.dataLength)
	for ; offset < readEnd; offset, dataOffset = offset+uint64(memBlockReadLength), dataOffset+memBlockReadLength {
		// prepare netBlock
		netBlockIndex = int32(offset / uint64(pNetINode.NetBlockCap))
		uNetBlock, err = p.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.netBlockDriver.ReleaseNetBlock(uNetBlock)

		// prepare memBlock
		memBlockIndex = int32(offset / uint64(pNetINode.MemBlockCap))
		memBlockStart = uint64(memBlockIndex) * uint64(pNetINode.MemBlockCap)
		memBlockReadOffset = int(offset - memBlockStart)
		if memBlockStart+uint64(pNetINode.MemBlockCap) < readEnd {
			// not the last block
			memBlockReadLength = int(memBlockStart + uint64(pNetINode.MemBlockCap) - offset)
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
		if arg.netQuery == nil {
			uMemBlock.Ptr().PReadWithMem(arg.data[dataOffset:dataOffset+memBlockReadLength], memBlockReadOffset)
		} else {
			// TODO read to connection
			err = uMemBlock.Ptr().PReadWithNetQuery(arg.netQuery, memBlockReadLength, memBlockReadOffset)
			if err != nil {
				goto READ_DATA_ONE_RUN_DONE
			}
		}

	READ_DATA_ONE_RUN_DONE:
		p.memBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
		if err != nil {
			goto READ_DATA_DONE
		}
	}

READ_DATA_DONE:
	return arg.dataLength, err
}

func (p *NetINodeDriver) PReadWithNetQuery(uNetINode solofsapitypes.NetINodeUintptr,
	netQuery *snettypes.NetQuery, dataLength int, offset uint64) (int, error) {
	return p.doPRead(uNetINode, preadArg{
		netQuery:   netQuery,
		data:       nil,
		dataLength: dataLength,
		offset:     offset,
	})
}

func (p *NetINodeDriver) PReadWithMem(uNetINode solofsapitypes.NetINodeUintptr,
	data []byte, offset uint64) (int, error) {
	return p.doPRead(uNetINode, preadArg{
		netQuery:   nil,
		data:       data,
		dataLength: len(data),
		offset:     offset,
	})
}
