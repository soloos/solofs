package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type MemBlockDriver struct {
	*soloosbase.SoloOSEnv
	tables map[int]*MemBlockTable
}

func (p *MemBlockDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	p.SoloOSEnv = soloOSEnv
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
func (p *MemBlockDriver) MustGetMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
	memBlockIndex int32) (solofsapitypes.MemBlockUintptr, bool) {
	var memBlockID soloosbase.PtrBindIndex
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) TryGetMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
	memBlockIndex int32) solofsapitypes.MemBlockUintptr {
	var memBlockID soloosbase.PtrBindIndex
	soloosbase.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].TryGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseMemBlockWithReadRelease(uMemBlock solofsapitypes.MemBlockUintptr) {
	if uMemBlock != 0 {
		p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseMemBlockWithReadRelease(uMemBlock)
	}
}

func (p *MemBlockDriver) MustGetTmpMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
	memBlockID soloosbase.PtrBindIndex) solofsapitypes.MemBlockUintptr {
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetTmpMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseTmpMemBlock(uMemBlock solofsapitypes.MemBlockUintptr) {
	p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseTmpMemBlock(uMemBlock)
}
