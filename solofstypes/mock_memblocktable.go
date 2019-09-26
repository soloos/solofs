package solofstypes

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
	"sync/atomic"
)

type MockMemBlockTable struct {
	*soloosbase.SoloOSEnv
	ichunkSize int
	mockID     int32
	hkvTable   offheap.HKVTableWithBytes12
}

func (p *MockMemBlockTable) Init(soloOSEnv *soloosbase.SoloOSEnv, ichunkSize int) error {
	var err error
	p.SoloOSEnv = soloOSEnv
	p.ichunkSize = ichunkSize
	err = p.OffheapDriver.InitHKVTableWithBytes12(&p.hkvTable, "mock",
		int(solofsapitypes.MemBlockStructSize+uintptr(p.ichunkSize)), -1, offheap.DefaultKVTableSharedCount,
		p.HKVTableInvokePrepareNewObject, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockTable) HKVTableInvokePrepareNewObject(uObject uintptr) {
	uMemBlock := solofsapitypes.MemBlockUintptr(uObject)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Bytes.Data = uObject + solofsapitypes.MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
	uMemBlock.Ptr().CompleteInit()
}

func (p *MockMemBlockTable) AllocMemBlock() solofsapitypes.MemBlockUintptr {
	var memBlockID soloosbase.PtrBindIndex
	id := atomic.AddInt32(&p.mockID, 1)
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(id), id)
	uObject, _ := p.hkvTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock := (solofsapitypes.MemBlockUintptr)(uObject)
	return uMemBlock
}
