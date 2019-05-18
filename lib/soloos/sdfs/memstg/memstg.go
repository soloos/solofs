package memstg

import (
	"soloos/common/sdfsapi"
	sdfsapitypes "soloos/common/sdfsapi/types"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/swalapi"
	"soloos/sdfs/netstg"
)

type MemStg struct {
	*soloosbase.SoloOSEnv

	NameNodeClient sdfsapi.NameNodeClient
	DataNodeClient sdfsapi.DataNodeClient
	netstg.NetBlockDriver
	MemBlockDriver
	NetINodeDriver

	swalClient swalapi.Client
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

	nameNodePeer = p.SNetDriver.AllocPeer(nameNodeSRPCServerAddr, sdfsapitypes.DefaultSDFSRPCProtocol)
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
