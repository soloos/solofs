package api

import (
	"soloos/sdfs/types"
	"soloos/snet"
)

type UploadMemBlockWithDisk func(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int) error

type DataNodeClient struct {
	snetClientDriver       *snet.ClientDriver
	uploadMemBlockWithDisk UploadMemBlockWithDisk
}

func (p *DataNodeClient) Init(snetClientDriver *snet.ClientDriver,
	uploadMemBlockWithDisk UploadMemBlockWithDisk) error {
	p.snetClientDriver = snetClientDriver
	p.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
	return nil
}

func (p *DataNodeClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}
