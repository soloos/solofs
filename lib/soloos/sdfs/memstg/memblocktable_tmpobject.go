package memstg

import (
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

func (p *MemBlockTable) MustGetTmpMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	var (
		uMemBlock types.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
	)
	uObject, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = types.MemBlockUintptr(uObject)
	return uMemBlock

}

func (p *MemBlockTable) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
