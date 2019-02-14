package memstg

import (
	"soloos/sdfs/types"
)

func (p *memBlockPoolHKVTable) beforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	if pMemBlock.IsInited() == false {
		return
	}
	pMemBlock.SetReleasable()
}

func (p *memBlockPoolHKVTable) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObjectByBytes12(uMemBlock.Ptr().ID)
}

func (p *memBlockPoolHKVTable) MustGetTmpMemBlockWithReadAcquire(blockID types.PtrBindIndex) types.MemBlockUintptr {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
	)
	u, _ = p.tmpMemBlockTable.MustGetObjectByBytes12WithReadAcquire(blockID)
	uObject = types.MemBlockUintptr(u)
	return uObject
}
