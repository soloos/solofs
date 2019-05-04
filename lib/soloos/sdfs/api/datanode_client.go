package api

import (
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
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
	*soloosbase.SoloOSEnv
	preadMemBlockWithDisk  PReadMemBlockWithDisk
	uploadMemBlockWithDisk UploadMemBlockWithDisk
}

func (p *DataNodeClient) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	p.SoloOSEnv = soloOSEnv
	return nil
}

func (p *DataNodeClient) SetPReadMemBlockWithDisk(preadMemBlockWithDisk PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *DataNodeClient) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}
