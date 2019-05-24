package types

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/sdbone/offheap"
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
		int(sdfsapitypes.MemBlockStructSize+uintptr(p.ichunkSize)), -1, offheap.DefaultKVTableSharedCount,
		p.HKVTableInvokePrepareNewObject, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockTable) HKVTableInvokePrepareNewObject(uObject uintptr) {
	uMemBlock := sdfsapitypes.MemBlockUintptr(uObject)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Bytes.Data = uObject + sdfsapitypes.MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
	uMemBlock.Ptr().CompleteInit()
}

func (p *MockMemBlockTable) AllocMemBlock() sdfsapitypes.MemBlockUintptr {
	var memBlockID soloosbase.PtrBindIndex
	id := atomic.AddInt32(&p.mockID, 1)
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(id), id)
	uObject, _ := p.hkvTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock := (sdfsapitypes.MemBlockUintptr)(uObject)
	return uMemBlock
}
