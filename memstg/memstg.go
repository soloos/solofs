package memstg

import (
	"soloos/common/solofsapi"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/solomqapi"
)

type MemStg struct {
	*soloosbase.SoloOSEnv

	SolonnClient solofsapi.SolonnClient
	SolodnClient solofsapi.SolodnClient
	NetBlockDriver
	MemBlockDriver
	NetINodeDriver

	solomqClient solomqapi.Client
}

func (p *MemStg) Init(soloOSEnv *soloosbase.SoloOSEnv,
	solonnPeer snettypes.Peer,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	var (
		err error
	)

	p.SoloOSEnv = soloOSEnv

	err = p.SolonnClient.Init(p.SoloOSEnv, solonnPeer.ID)
	if err != nil {
		return err
	}

	err = p.SolodnClient.Init(p.SoloOSEnv)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloOSEnv,
		&p.SolonnClient, &p.SolodnClient,
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
		&p.SolonnClient,
		p.NetINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.NetINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.NetINodeDriver.NetINodeCommitSizeInDB,
	)
	if err != nil {
		return err
	}

	return nil
}
