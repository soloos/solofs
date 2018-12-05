package netstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"sync"
)

type NetBlockPool struct {
	options NetBlockPoolOptions
	driver  *NetBlockDriver

	netBlockAllocRWMutex sync.RWMutex
	offheapPool          offheap.RawObjectPool
	pool                 map[types.INodeBlockID]types.NetBlockUintptr
}

func (p *NetBlockPool) Init(options NetBlockPoolOptions,
	driver *NetBlockDriver) error {
	var err error

	p.options = options

	p.driver = driver

	err = p.driver.offheapDriver.InitRawObjectPool(&p.offheapPool,
		int(types.NetBlockStructSize), p.options.RawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[types.INodeBlockID]types.NetBlockUintptr)

	return nil
}

func (p *NetBlockPool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

func (p *NetBlockPool) PrepareNetBlockMetadata(uINode types.INodeUintptr, uNetBlock types.NetBlockUintptr) {
	pNetBlock := uNetBlock.Ptr()
	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	// p.driver.snetClientDriver.Call()

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockPool) MustGetBlock(uINode types.INodeUintptr,
	netBlockIndex int) types.NetBlockUintptr {
	var (
		netBlockID types.INodeBlockID
		uNetBlock  types.NetBlockUintptr
		exists     bool
	)
	types.EncodeINodeBlockID(&netBlockID, uINode.Ptr().ID, netBlockIndex)
	p.netBlockAllocRWMutex.RLock()
	uNetBlock, exists = p.pool[netBlockID]
	p.netBlockAllocRWMutex.RUnlock()

	if !exists {
		p.netBlockAllocRWMutex.Lock()
		uNetBlock, exists = p.pool[netBlockID]
		if !exists {
			uNetBlock = types.NetBlockUintptr(p.offheapPool.AllocRawObject())
			p.pool[netBlockID] = uNetBlock
		}
		p.netBlockAllocRWMutex.Unlock()

		if !uNetBlock.Ptr().IsMetaDataInited {
			p.PrepareNetBlockMetadata(uINode, uNetBlock)
		}
	}

	return uNetBlock
}
