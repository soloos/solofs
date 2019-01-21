package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
	"sync"
)

type PrepareNetBlockMetaData func(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int) error

type NetBlockDriverHelper struct {
	*api.NameNodeClient
	PrepareNetBlockMetaData
}

type NetBlockDriver struct {
	helper NetBlockDriverHelper

	offheapDriver        *offheap.OffheapDriver
	netBlockAllocRWMutex sync.RWMutex
	offheapPool          offheap.RawObjectPool
	netBlockPool         types.NetBlockPool

	snetDriver       *snet.NetDriver
	snetClientDriver *snet.ClientDriver
	dataNodeClient   *api.DataNodeClient

	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.netBlockPool.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.dataNodeClient = dataNodeClient

	p.SetHelper(nameNodeClient, prepareNetBlockMetaData)

	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) SetHelper(
	nameNodeClient *api.NameNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
}

func (p *NetBlockDriver) SetPReadMemBlockWithDisk(preadWithDisk api.PReadMemBlockWithDisk) {
	p.dataNodeClient.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *NetBlockDriver) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk api.UploadMemBlockWithDisk) {
	p.dataNodeClient.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}

func (p *NetBlockDriver) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockDriver) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode types.NetINodeUintptr,
	netBlockIndex int) (types.NetBlockUintptr, error) {
	var (
		uNetBlock types.NetBlockUintptr
		pNetBlock *types.NetBlock
		isLoaded  bool
		err       error
	)

	uNetBlock, isLoaded = p.netBlockPool.MustGetNetBlock(uNetINode, netBlockIndex)
	pNetBlock = uNetBlock.Ptr()
	if isLoaded == false || uNetBlock.Ptr().IsDBMetaDataInited == false {
		pNetBlock.DBMetaDataInitMutex.Lock()
		if pNetBlock.IsDBMetaDataInited == false {
			err = p.helper.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
		}
		pNetBlock.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		pNetBlock.SharedPointer.SetReleasable()
		p.netBlockPool.ReleaseNetBlock(uNetBlock)
		return 0, err
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) FlushMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr) error {
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()
	return nil
}

func (p *NetBlockDriver) GetDataNodeClient() *api.DataNodeClient {
	return p.dataNodeClient
}
