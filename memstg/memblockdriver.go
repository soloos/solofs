package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
)

type MemBlockDriver struct {
	*soloosbase.SoloosEnv
	tables map[int]*MemBlockTable
}

func (p *MemBlockDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	p.SoloosEnv = soloosEnv
	p.tables = make(map[int]*MemBlockTable)
	for _, memBlockTableOptions := range memBlockDriverOptions.MemBlockTableOptionsList {
		p.PrepareMemBlockTable(memBlockTableOptions)
	}
	return nil
}

func (p *MemBlockDriver) PrepareMemBlockTable(memBlockTableOptions MemBlockTableOptions) error {
	if _, exists := p.tables[memBlockTableOptions.ObjectSize]; exists {
		return nil
	}

	var (
		memBlockTable *MemBlockTable
		err           error
	)

	memBlockTable = new(MemBlockTable)
	err = memBlockTable.Init(memBlockTableOptions, p)
	if err != nil {
		return err
	}

	p.tables[memBlockTable.options.ObjectSize] = memBlockTable
	return nil
}

// MustGetMemBlockWithReadAcquire get or init a memblock's offheap
func (p *MemBlockDriver) MustGetMemBlockWithReadAcquire(uNetINode solofstypes.NetINodeUintptr,
	memBlockIndex int32) (solofstypes.MemBlockUintptr, bool) {
	var memBlockID soloosbase.PtrBindIndex
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) TryGetMemBlockWithReadAcquire(uNetINode solofstypes.NetINodeUintptr,
	memBlockIndex int32) solofstypes.MemBlockUintptr {
	var memBlockID soloosbase.PtrBindIndex
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].TryGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseMemBlockWithReadRelease(uMemBlock solofstypes.MemBlockUintptr) {
	if uMemBlock != 0 {
		p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseMemBlockWithReadRelease(uMemBlock)
	}
}

func (p *MemBlockDriver) MustGetTmpMemBlockWithReadAcquire(uNetINode solofstypes.NetINodeUintptr,
	memBlockID soloosbase.PtrBindIndex) solofstypes.MemBlockUintptr {
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetTmpMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseTmpMemBlock(uMemBlock solofstypes.MemBlockUintptr) {
	p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseTmpMemBlock(uMemBlock)
}
