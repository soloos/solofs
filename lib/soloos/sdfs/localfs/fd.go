package localfs

import (
	"os"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"sync/atomic"
)

type Fd struct {
	accessor  int32
	uNetINode types.NetINodeUintptr
	filePath  string
	file      *os.File
}

func (p *Fd) Init(uNetINode types.NetINodeUintptr, filePath string) error {
	var err error
	p.uNetINode = uNetINode
	p.filePath = filePath
	p.file, err = os.OpenFile(p.filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (p *Fd) Upload(uJob types.UploadMemBlockJobUintptr) error {
	var (
		req                 snettypes.Request
		netINodeWriteOffset int64
		memBlockCap         int
		pChunkMask          *offheap.ChunkMask
		writeData           []byte
		err                 error
	)

	pChunkMask = uJob.Ptr().UploadMaskProcessing.Ptr()

	req.OffheapBody.OffheapBytes = uJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockCap = uJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		netINodeWriteOffset = int64(memBlockCap)*int64(uJob.Ptr().MemBlockIndex) +
			int64(pChunkMask.MaskArray[chunkMaskIndex].Offset)

		writeData = (*uJob.Ptr().UMemBlock.Ptr().BytesSlice())[pChunkMask.MaskArray[chunkMaskIndex].Offset:pChunkMask.MaskArray[chunkMaskIndex].End]
		err = p.WriteAt(writeData, netINodeWriteOffset)
		if err != nil {
			goto PWRITE_DONE
		}
	}

PWRITE_DONE:
	return err
}

func (p *Fd) WriteAt(data []byte, netINodeOffset int64) error {
	var (
		off int
		n   int
		err error
	)
	for off = 0; off < len(data); off += n {
		n, err = p.file.WriteAt(data, netINodeOffset+int64(off))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Fd) BorrowResource() {
	atomic.AddInt32(&p.accessor, 1)
}

func (p *Fd) ReturnResource() {
	atomic.AddInt32(&p.accessor, -1)
}
