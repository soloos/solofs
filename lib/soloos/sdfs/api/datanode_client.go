package api

import (
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

type PReadMemBlockWithDisk func(uNetINode types.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int32,
	uMemBlock types.MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error)
type UploadMemBlockWithDisk func(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int) error

type DataNodeClient struct {
	snetClientDriver       *snet.ClientDriver
	preadMemBlockWithDisk  PReadMemBlockWithDisk
	uploadMemBlockWithDisk UploadMemBlockWithDisk
}

func (p *DataNodeClient) Init(snetClientDriver *snet.ClientDriver) error {
	p.snetClientDriver = snetClientDriver
	return nil
}

func (p *DataNodeClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *DataNodeClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}
