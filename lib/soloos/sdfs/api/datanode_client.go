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
	offset int64, length int) error
type UploadMemBlockWithDisk func(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int) error

type DataNodeClient struct {
	snetClientDriver       *snet.ClientDriver
	preadMemBlockWithDisk  PReadMemBlockWithDisk
	uploadMemBlockWithDisk UploadMemBlockWithDisk
}

func (p *DataNodeClient) Init(snetClientDriver *snet.ClientDriver,
	preadWithDisk PReadMemBlockWithDisk,
	uploadMemBlockWithDisk UploadMemBlockWithDisk) error {
	p.snetClientDriver = snetClientDriver
	p.SetPReadMemBlockWithDisk(preadWithDisk)
	p.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
	return nil
}

func (p *DataNodeClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *DataNodeClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}
