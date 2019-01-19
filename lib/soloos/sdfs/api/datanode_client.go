package api

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
)

type PReadMemBlockWithDisk func(uNetINode types.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int,
	uMemBlock types.MemBlockUintptr, memBlockIndex int,
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
