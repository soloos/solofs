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

type NetBlockDriver struct {
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

func (p *NetBlockDriver) Init(offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	err = p.netBlockPool.Init(-1, p.offheapDriver)
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
	if isLoaded == false || uNetBlock.Ptr().IsMetaDataInited == false {
		pNetBlock.MetaDataInitMutex.Lock()
		if pNetBlock.IsMetaDataInited == false {
			err = p.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
			if err == nil {
				pNetBlock.IsMetaDataInited = true
			}
		}
		pNetBlock.MetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetBlock
		return 0, err
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo protocol.NetINodeNetBlockInfoResponse
		backend      protocol.NetBlockBackend
		peerID       snettypes.PeerID
		uPeer        snettypes.PeerUintptr
		i            int
		err          error
	)

	err = p.nameNodeClient.PrepareNetBlockMetaData(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
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
	pNetBlock.SyncDataPrimaryBackendTransferCount = pNetBlock.SyncDataBackends.Len - 1

	pNetBlock.IsSyncDataBackendsInited = true

	return nil
}

func (p *NetBlockDriver) FlushMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr) error {
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()
	return nil
}
