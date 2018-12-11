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
	netBlockPool         types.NetBlockPool

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

	err = p.netBlockPool.Init(p.options.RawChunksLimit, p.offheapDriver)
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
func (p *NetBlockDriver) MustGetBlock(uNetINode types.NetINodeUintptr,
	netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		exists    bool
		err       error
	)

	uNetBlock, exists = p.netBlockPool.MustGetNetBlock(uNetINode, netBlockIndex)

	if exists == false || uNetBlock.Ptr().IsMetaDataInited == false {
		err = p.prepareNetBlockMetadata(uNetINode, netBlockIndex, uNetBlock)
		if err != nil {
			return 0, err
		}
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) prepareNetBlockMetadata(uNetINode types.NetINodeUintptr,
	netblockIndex int,
	uNetBlock types.NetBlockUintptr,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.nameNodeClient.PrepareNetBlockMetadata(uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		goto PREPARE_DONE
	}

	pNetBlock.IsMetaDataInited = true

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
	return err
}
