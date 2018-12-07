package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
	"sync"
)

type NetBlockDriver struct {
	options              NetBlockDriverOptions
	offheapDriver        *offheap.OffheapDriver
	netBlockAllocRWMutex sync.RWMutex
	offheapPool          offheap.RawObjectPool
	pool                 map[types.INodeBlockID]types.NetBlockUintptr

	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver
	nameNodeClient   *api.NameNodeClient
	dataNodeClient   *api.DataNodeClient

	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(options NetBlockDriverOptions,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) error {
	var err error
	p.options = options
	p.offheapDriver = offheapDriver
	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	err = p.offheapDriver.InitRawObjectPool(&p.offheapPool,
		int(types.NetBlockStructSize), p.options.RawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	p.pool = make(map[types.INodeBlockID]types.NetBlockUintptr)

	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.nameNodeClient = nameNodeClient
	p.dataNodeClient = dataNodeClient

	return nil
}

func (p *NetBlockDriver) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockDriver) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetBlock(uINode types.INodeUintptr,
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

	if exists == false {
		p.netBlockAllocRWMutex.Lock()
		uNetBlock, exists = p.pool[netBlockID]
		if exists == false {
			uNetBlock = types.NetBlockUintptr(p.offheapPool.AllocRawObject())
			p.pool[netBlockID] = uNetBlock
		}
		p.netBlockAllocRWMutex.Unlock()
	}

	if uNetBlock.Ptr().IsMetaDataInited == false {
		p.PrepareNetBlockMetadata(uINode, netBlockIndex, uNetBlock)
	}

	return uNetBlock
}
