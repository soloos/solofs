package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := solofstypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

func (p *MemBlockTable) MustGetTmpMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) solofstypes.MemBlockUintptr {
	var (
		uMemBlock solofstypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
	)
	uObject, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = solofstypes.MemBlockUintptr(uObject)
	return uMemBlock

}

func (p *MemBlockTable) ReleaseTmpMemBlock(uMemBlock solofstypes.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
