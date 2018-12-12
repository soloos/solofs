package api

import (
	"soloos/snet"
	snettypes "soloos/snet/types"
)

type NameNodeClient struct {
	snetClientDriver *snet.ClientDriver
	nameNodePeer     snettypes.PeerUintptr
}

func (p *NameNodeClient) Init(snetClientDriver *snet.ClientDriver,
	nameNodePeer snettypes.PeerUintptr) error {
	p.snetClientDriver = snetClientDriver
	p.nameNodePeer = nameNodePeer
	return nil
}
