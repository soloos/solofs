package memstg

import (
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) PWrite(uNetINode types.NetINodeUintptr, data []byte, offset int64) error {
	var (
		isSuccess           bool
		memBlockIndex       int
		memBlockBytesOffset int
		memBlockBytesEnd    int
		uMemBlock           types.MemBlockUintptr
		netBlockIndex       int
		uNetBlock           types.NetBlockUintptr
		writeEnd            int64
		pNetINode           = uNetINode.Ptr()
		i                   int
		err                 error
	)

	pNetINode.WriteDataRWMutex.RLock()

	for writeEnd = offset + int64(len(data)); offset < writeEnd; offset += int64(pNetINode.MemBlockCap) {
		// prepare netBlock
		netBlockIndex = int(offset / int64(pNetINode.NetBlockCap))
		uNetBlock, err = p.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)

		// prepare memBlock
		memBlockIndex = int(offset / int64(pNetINode.MemBlockCap))
		memBlockBytesOffset = int(offset - int64(memBlockIndex)*int64(pNetINode.MemBlockCap))
		memBlockBytesEnd = memBlockBytesOffset + len(data)
		uMemBlock, _ = p.memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)

		// write in memblock
		for i = 0; i < 6; i++ {
			isSuccess = uMemBlock.Ptr().PWrite(data, memBlockBytesOffset)
			if isSuccess == false {
				err = p.unsafeMemBlockRebaseNetBlock(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex)
				if err != nil {
					goto WRITE_DATA_ONE_RUN_DONE
				}
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
			memBlockBytesOffset, memBlockBytesEnd)
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

func (p *NetINodeDriver) Flush(uNetINode types.NetINodeUintptr) error {
	// TODO common offset in metadb
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
