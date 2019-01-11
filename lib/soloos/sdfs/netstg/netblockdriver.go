package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"sync"
)

type PrepareNetBlockMetaData func(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int) error

type NetBlockDriverHelper struct {
	nameNodeClient          *api.NameNodeClient
	PrepareNetBlockMetaData PrepareNetBlockMetaData
}

type NetBlockDriver struct {
	Helper NetBlockDriverHelper

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
	err = p.netBlockPool.Init(-1, p.offheapDriver)
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
	p.Helper.nameNodeClient = nameNodeClient
	p.Helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
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
			err = p.Helper.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
		}
		pNetBlock.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetBlock
		return 0, err
	}

	return uNetBlock, nil
}

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int,
	syncDataPrimaryBackendTransferCount int,
) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo protocol.NetINodeNetBlockInfoResponse
		backend      protocol.SNetPeer
		peerID       snettypes.PeerID
		uPeer        snettypes.PeerUintptr
		i            int
		err          error
	)

	err = p.Helper.nameNodeClient.PrepareNetBlockMetaData(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		return err
	}

	pNetBlock.StorDataBackends.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		netBlockInfo.Backends(&backend, i)
		copy(peerID[:], backend.PeerID())
		uPeer, _ = p.snetDriver.MustGetPeer(&peerID, string(backend.Address()), types.DefaultSDFSRPCProtocol)
		pNetBlock.StorDataBackends.Append(uPeer)
	}

	pNetBlock.SyncDataBackends = pNetBlock.StorDataBackends
	pNetBlock.SyncDataPrimaryBackendTransferCount = syncDataPrimaryBackendTransferCount
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaDataWithTransfer(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int) error {
	var err error
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex,
		uNetBlock.Ptr().SyncDataBackends.Len-1)
	if err != nil {
		return err
	}
	uNetBlock.Ptr().IsDBMetaDataInited = true
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaDataWithFanout(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int) error {
	var err error
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex, 0)
	if err != nil {
		return err
	}
	uNetBlock.Ptr().IsDBMetaDataInited = true
	return nil
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
