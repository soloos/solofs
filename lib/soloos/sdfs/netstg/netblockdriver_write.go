package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) WriteBytesAt(uINode types.INodeUintptr, uNetBlock types.NetBlockUintptr,
	bytes []byte, offset int64) error {
	var err error
	pNetBlock := uNetBlock.Ptr()
	for i := 0; i < pNetBlock.DataNodes.Len; i++ {
		err = p.snetClientDriver.Write(pNetBlock.DataNodes.Arr[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *NetBlockDriver) WriteMemBlockAt(uINode types.INodeUintptr, uNetBlock types.NetBlockUintptr,
	memBlock types.MemBlockUintptr,
	memBlockDataOffset, memBlockDataLen int,
	offset int64) error {
	var bytes = (*memBlock.Ptr().BytesSlice())[memBlockDataOffset : memBlockDataOffset+memBlockDataLen]
	return p.WriteBytesAt(uINode, uNetBlock, bytes, offset)
}
