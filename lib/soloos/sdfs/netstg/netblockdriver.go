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
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo protocol.NetINodeNetBlockInfoResponse
		backend      protocol.NetBlockBackend
		peerID       snettypes.PeerID
		uPeer        snettypes.PeerUintptr
		i            int
		err          error
	)

	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.nameNodeClient.PrepareNetBlockMetadata(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		goto PREPARE_DONE
	}

	pNetBlock.DataNodes.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		netBlockInfo.Backends(&backend, i)
		copy(peerID[:], backend.PeerID())
		uPeer, _ = p.snetDriver.MustGetPeer(&peerID, string(backend.Address()), types.DefaultSDFSRPCProtocol)
		pNetBlock.DataNodes.Append(uPeer)
	}

	pNetBlock.IsMetaDataInited = true

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
	return err
}
