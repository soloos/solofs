package memstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/sdbone/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := sdfsapitypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

func (p *MemBlockTable) MustGetTmpMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) sdfsapitypes.MemBlockUintptr {
	var (
		uMemBlock sdfsapitypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
	)
	uObject, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = sdfsapitypes.MemBlockUintptr(uObject)
	return uMemBlock

}

func (p *MemBlockTable) ReleaseTmpMemBlock(uMemBlock sdfsapitypes.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
