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
	pool                 types.NetBlockPool

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

	err = p.pool.Init(p.options.RawChunksLimit, p.offheapDriver)
	if err != nil {
		return err
	}

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
	netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		exists    bool
		err       error
	)

	uNetBlock, exists = p.pool.MustGetNetBlock(uINode, netBlockIndex)

	if exists == false || uNetBlock.Ptr().IsMetaDataInited == false {
		err = p.PrepareNetBlockMetadata(uINode, netBlockIndex, uNetBlock)
		if err != nil {
			return 0, err
		}
	}

	return uNetBlock, nil
}
