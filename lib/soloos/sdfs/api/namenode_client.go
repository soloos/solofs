package api

import (
	"soloos/snet"
	snettypes "soloos/snet/types"
)

type NameNodeClient struct {
	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver
	nameNodePeer     snettypes.PeerUintptr
}

func (p *NameNodeClient) Init(snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodePeer snettypes.PeerUintptr) error {
	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.nameNodePeer = nameNodePeer
	return nil
}
