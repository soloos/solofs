package memstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/sdbone/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := sdfsapitypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockTable) MustGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) (sdfsapitypes.MemBlockUintptr, bool) {
	var (
		uMemBlock sdfsapitypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
		loaded    bool
	)
	uObject, loaded = p.memBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = sdfsapitypes.MemBlockUintptr(uObject)
	return uMemBlock, loaded
}

func (p *MemBlockTable) TryGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) sdfsapitypes.MemBlockUintptr {
	var uObject sdfsapitypes.MemBlockUintptr
	uObject = sdfsapitypes.MemBlockUintptr(p.memBlockTable.TryGetObjectWithReadAcquire(memBlockID))
	return uObject
}

func (p *MemBlockTable) ReleaseMemBlockWithReadRelease(uMemBlock sdfsapitypes.MemBlockUintptr) {
	p.memBlockTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithBytes12(uMemBlock))
}
