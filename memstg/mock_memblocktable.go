package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
	"sync/atomic"
)

type MockMemBlockTable struct {
	*soloosbase.SoloosEnv
	ichunkSize int
	mockID     int32
	hkvTable   offheap.HKVTableWithBytes12
}

func (p *MockMemBlockTable) Init(soloosEnv *soloosbase.SoloosEnv, ichunkSize int) error {
	var err error
	p.SoloosEnv = soloosEnv
	p.ichunkSize = ichunkSize
	err = p.OffheapDriver.InitHKVTableWithBytes12(&p.hkvTable, "mock",
		int(solofstypes.MemBlockStructSize+uintptr(p.ichunkSize)), -1, offheap.DefaultKVTableSharedCount,
		p.HKVTableInvokePrepareNewObject, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockTable) HKVTableInvokePrepareNewObject(uObject uintptr) {
	uMemBlock := solofstypes.MemBlockUintptr(uObject)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Bytes.Data = uObject + solofstypes.MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
	uMemBlock.Ptr().CompleteInit()
}

func (p *MockMemBlockTable) AllocMemBlock() solofstypes.MemBlockUintptr {
	var memBlockID soloosbase.PtrBindIndex
	id := atomic.AddInt32(&p.mockID, 1)
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(id), id)
	uObject, _ := p.hkvTable.MustGetObjectWithReadAcquire(memBlockID)
	uMemBlock := (solofstypes.MemBlockUintptr)(uObject)
	return uMemBlock
}
