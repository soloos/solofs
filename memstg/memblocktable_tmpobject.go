package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := solofsapitypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

func (p *MemBlockTable) MustGetTmpMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) solofsapitypes.MemBlockUintptr {
	var (
		uMemBlock solofsapitypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
	)
	uObject, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = solofsapitypes.MemBlockUintptr(uObject)
	return uMemBlock

}

func (p *MemBlockTable) ReleaseTmpMemBlock(uMemBlock solofsapitypes.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
