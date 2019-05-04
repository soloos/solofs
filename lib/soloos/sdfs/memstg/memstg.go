package memstg

import (
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
)

type MemStg struct {
	*soloosbase.SoloOSEnv

	NameNodeClient api.NameNodeClient
	DataNodeClient api.DataNodeClient
	netstg.NetBlockDriver
	MemBlockDriver
	NetINodeDriver
}

func (p *MemStg) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	var (
		nameNodePeer snettypes.PeerUintptr
		err          error
	)

	p.SoloOSEnv = soloOSEnv

	nameNodePeer = p.SNetDriver.AllocPeer(nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	err = p.NameNodeClient.Init(p.SoloOSEnv, nameNodePeer)
	if err != nil {
		return err
	}

	err = p.DataNodeClient.Init(p.SoloOSEnv)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloOSEnv,
		&p.NameNodeClient, &p.DataNodeClient,
		p.NetBlockDriver.PrepareNetBlockMetaDataWithFanout,
	)
	if err != nil {
		return err
	}

	err = p.MemBlockDriver.Init(p.SoloOSEnv, memBlockDriverOptions)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.SoloOSEnv, &p.NetBlockDriver, &p.MemBlockDriver,
		&p.NameNodeClient,
		p.NetINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.NetINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.NetINodeDriver.NetINodeCommitSizeInDB,
	)
	if err != nil {
		return err
	}

	return nil
}
