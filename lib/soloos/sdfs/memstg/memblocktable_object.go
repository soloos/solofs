package memstg

import (
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockTable) MustGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
		loaded  bool
	)
	u, loaded = p.memBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uObject = types.MemBlockUintptr(u)
	return uObject, loaded
}

func (p *MemBlockTable) TryGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	var uObject types.MemBlockUintptr
	uObject = types.MemBlockUintptr(p.memBlockTable.TryGetObjectWithReadAcquire(memBlockID))
	return uObject
}

func (p *MemBlockTable) ReleaseMemBlockWithReadRelease(uMemBlock types.MemBlockUintptr) {
	p.memBlockTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithBytes12(uMemBlock))
}
