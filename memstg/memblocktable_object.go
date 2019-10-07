package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

func (p *MemBlockTable) hkvTableInvokeBeforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := solofstypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.Reset()
	pMemBlock.SetReleasable()
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockTable) MustGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) (solofstypes.MemBlockUintptr, bool) {
	var (
		uMemBlock solofstypes.MemBlockUintptr
		uObject   offheap.HKVTableObjectUPtrWithBytes12
		loaded    bool
	)
	uObject, loaded = p.memBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock = solofstypes.MemBlockUintptr(uObject)
	return uMemBlock, loaded
}

func (p *MemBlockTable) TryGetMemBlockWithReadAcquire(memBlockID soloosbase.PtrBindIndex) solofstypes.MemBlockUintptr {
	var uObject solofstypes.MemBlockUintptr
	uObject = solofstypes.MemBlockUintptr(p.memBlockTable.TryGetObjectWithReadAcquire(memBlockID))
	return uObject
}

func (p *MemBlockTable) ReleaseMemBlockWithReadRelease(uMemBlock solofstypes.MemBlockUintptr) {
	p.memBlockTable.ReadReleaseObject(offheap.HKVTableObjectUPtrWithBytes12(uMemBlock))
}
