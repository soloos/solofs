package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := solofsapitypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockTable) MustGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) (solofsapitypes.MemBlockUintptr, bool) {
	var (
		uMemBlock solofsapitypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
		loaded    bool
	)
	uObject, loaded = p.memBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = solofsapitypes.MemBlockUintptr(uObject)
	return uMemBlock, loaded
}

func (p *MemBlockTable) TryGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) solofsapitypes.MemBlockUintptr {
	var uObject solofsapitypes.MemBlockUintptr
	uObject = solofsapitypes.MemBlockUintptr(p.memBlockTable.TryGetObjectWithReadAcquire(memBlockID))
	return uObject
}

func (p *MemBlockTable) ReleaseMemBlockWithReadRelease(uMemBlock solofsapitypes.MemBlockUintptr) {
	p.memBlockTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithBytes12(uMemBlock))
}
