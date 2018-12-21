package memstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

type pwriteArg struct {
	conn       *snettypes.Connection
	dataLength int
	data       []byte
	offset     int64
}

func (p *NetINodeDriver) doPWrite(uNetINode types.NetINodeUintptr,
	arg pwriteArg) error {
	var (
		isSuccess           bool
		memBlockIndex       int
		memBlockStart       int64
		memBlockWriteOffset int
		memBlockWriteEnd    int
		memBlockWriteLength int
		uMemBlock           types.MemBlockUintptr
		netBlockIndex       int
		uNetBlock           types.NetBlockUintptr
		offset              = arg.offset
		dataOffset          = 0
		writeEnd            int64
		pNetINode           = uNetINode.Ptr()
		i                   int
		err                 error
	)

	pNetINode.WriteDataRWMutex.RLock()

	writeEnd = offset + int64(arg.dataLength)
	for ; offset < writeEnd; offset, dataOffset = offset+int64(pNetINode.MemBlockCap), dataOffset+pNetINode.MemBlockCap {
		// prepare netBlock
		netBlockIndex = int(offset / int64(pNetINode.NetBlockCap))
		uNetBlock, err = p.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)

		// prepare memBlock
		memBlockIndex = int(offset / int64(pNetINode.MemBlockCap))
		memBlockStart = int64(memBlockIndex) * int64(pNetINode.MemBlockCap)
		memBlockWriteOffset = int(offset - memBlockStart)
		if memBlockStart+int64(pNetINode.MemBlockCap) < writeEnd {
			// not the last block
			memBlockWriteLength = int(memBlockStart + int64(pNetINode.MemBlockCap) - offset)
		} else {
			// the last block
			memBlockWriteLength = int(writeEnd - offset)
		}
		memBlockWriteEnd = memBlockWriteOffset + memBlockWriteLength
		uMemBlock, _ = p.memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)

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
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
		if err != nil {
			goto WRITE_DATA_DONE
		}
	}

WRITE_DATA_DONE:
	pNetINode.WriteDataRWMutex.RUnlock()
	return err
}

func (p *NetINodeDriver) PWriteWithConn(uNetINode types.NetINodeUintptr,
	conn *snettypes.Connection, dataLength int, offset int64) error {
	return p.doPWrite(uNetINode, pwriteArg{
		conn:       conn,
		data:       nil,
		dataLength: dataLength,
		offset:     offset,
	})
}

func (p *NetINodeDriver) PWriteWithMem(uNetINode types.NetINodeUintptr,
	data []byte, offset int64) error {
	return p.doPWrite(uNetINode, pwriteArg{
		conn:       nil,
		data:       data,
		dataLength: len(data),
		offset:     offset,
	})
}

func (p *NetINodeDriver) Flush(uNetINode types.NetINodeUintptr) error {
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
	return err
}
