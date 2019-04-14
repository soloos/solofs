package memstg

import "soloos/sdfs/types"

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

func (p *MemBlockTable) MustGetTmpMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
	)
	u, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uObject = types.MemBlockUintptr(u)
	return uObject

}

func (p *MemBlockTable) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
