package memstg

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

type pwriteArg struct {
	conn       *snettypes.Connection
	dataLength int
	data       []byte
	offset     uint64
}

func (p *NetINodeDriver) doPWrite(uNetINode types.NetINodeUintptr,
	arg pwriteArg) error {
	var (
		isSuccess           bool
		uMemBlock           types.MemBlockUintptr
		uNetBlock           types.NetBlockUintptr
		memBlockIndex       int32
		netBlockIndex       int32
		memBlockStart       uint64
		memBlockWriteOffset int
		memBlockWriteEnd    int
		memBlockWriteLength int
		offset              = arg.offset
		dataOffset          = 0
		writeEnd            uint64
		pNetINode           = uNetINode.Ptr()
		i                   int
		err                 error
	)

	pNetINode.WriteDataRWMutex.RLock()

	writeEnd = offset + uint64(arg.dataLength)
	for ; offset < writeEnd; offset, dataOffset = offset+uint64(memBlockWriteLength), dataOffset+memBlockWriteLength {
		// prepare netBlock
		netBlockIndex = int32(offset / uint64(pNetINode.NetBlockCap))
		uNetBlock, err = p.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.netBlockDriver.ReleaseNetBlock(uNetBlock)

		// prepare memBlock
		memBlockIndex = int32(offset / uint64(pNetINode.MemBlockCap))
		memBlockStart = uint64(memBlockIndex) * uint64(pNetINode.MemBlockCap)
		memBlockWriteOffset = int(offset - memBlockStart)
		if memBlockStart+uint64(pNetINode.MemBlockCap) < writeEnd {
			// not the last block
			memBlockWriteLength = int(memBlockStart + uint64(pNetINode.MemBlockCap) - offset)
		} else {
			// the last block
			memBlockWriteLength = int(writeEnd - offset)
		}
		memBlockWriteEnd = memBlockWriteOffset + memBlockWriteLength
		uMemBlock, _ = p.memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)

		// TODO refine me
		// write in memblock
		for i = 0; i < 6; i++ {
			if arg.conn == nil {
				isSuccess = uMemBlock.Ptr().PWriteWithMem(arg.data[dataOffset:dataOffset+memBlockWriteLength],
					memBlockWriteOffset)
			} else {
				isSuccess = uMemBlock.Ptr().PWriteWithConn(arg.conn, memBlockWriteLength, memBlockWriteOffset)
			}
			if isSuccess {
				break
			}
			err = p.unsafeMemBlockRebaseNetBlock(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex)
			if err != nil {
				goto WRITE_DATA_ONE_RUN_DONE
			}
		}
		if isSuccess == false {
			// TODO catch error
			err = types.ErrRetryTooManyTimes
			goto WRITE_DATA_ONE_RUN_DONE
		}

		// write in netblock
		if err != nil {
			goto WRITE_DATA_ONE_RUN_DONE
		}

		err = p.netBlockDriver.PWrite(uNetINode,
			uNetBlock, netBlockIndex,
			uMemBlock, memBlockIndex,
			memBlockWriteOffset, memBlockWriteEnd)
		if err != nil {
			goto WRITE_DATA_ONE_RUN_DONE
		}

	WRITE_DATA_ONE_RUN_DONE:
		p.memBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
		if err != nil {
			goto WRITE_DATA_DONE
		}
	}

WRITE_DATA_DONE:
	if pNetINode.Size < writeEnd {
		pNetINode.Size = writeEnd
	}

	pNetINode.WriteDataRWMutex.RUnlock()
	return err
}

func (p *NetINodeDriver) PWriteWithConn(uNetINode types.NetINodeUintptr,
	conn *snettypes.Connection, dataLength int, offset uint64) error {
	return p.doPWrite(uNetINode, pwriteArg{
		conn:       conn,
		data:       nil,
		dataLength: dataLength,
		offset:     offset,
	})
}

func (p *NetINodeDriver) PWriteWithMem(uNetINode types.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.doPWrite(uNetINode, pwriteArg{
		conn:       nil,
		data:       data,
		dataLength: len(data),
		offset:     offset,
	})
}

func (p *NetINodeDriver) Sync(uNetINode types.NetINodeUintptr) error {
	// TODO commit offset in metadb
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)
	pNetINode.WriteDataRWMutex.Lock()
	pNetINode.SyncDataSig.Wait()
	pNetINode.WriteDataRWMutex.Unlock()
	err = pNetINode.LastSyncDataError
	pNetINode.LastSyncDataError = nil

	if pNetINode.LastCommitSize == pNetINode.Size {
		return nil
	}

	// TODO improve me
	err = p.NetINodeCommitSizeInDB(uNetINode, pNetINode.Size)
	if err != nil {
		return err
	}

	return nil
}
