package api

import (
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
)

type NameNodeClient struct {
	*soloosbase.SoloOSEnv
	nameNodePeer snettypes.PeerUintptr
}

func (p *NameNodeClient) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodePeer snettypes.PeerUintptr) error {
	p.SoloOSEnv = soloOSEnv
	p.nameNodePeer = nameNodePeer
	return nil
}
