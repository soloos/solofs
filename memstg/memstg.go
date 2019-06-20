package memstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
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
	nameNodePeer snettypes.Peer,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	var (
		err error
	)

	p.SoloOSEnv = soloOSEnv

	err = p.NameNodeClient.Init(p.SoloOSEnv, nameNodePeer.ID)
	if err != nil {
		return err
	}

	err = p.DataNodeClient.Init(p.SoloOSEnv)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloOSEnv,
		&p.NameNodeClient, &p.DataNodeClient,
		p.NetBlockDriver.PrepareNetBlockMetaData,
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
