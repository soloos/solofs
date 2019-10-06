package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapi"
	"soloos/common/solomqapi"
	"soloos/common/soloosbase"
)

type MemStg struct {
	*soloosbase.SoloosEnv

	SolonnClient solofsapi.SolonnClient
	SolodnClient solofsapi.SolodnClient
	NetBlockDriver
	MemBlockDriver
	NetINodeDriver

	solomqClient solomqapi.Client
}

func (p *MemStg) Init(soloosEnv *soloosbase.SoloosEnv,
	solonnPeer snettypes.Peer,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	var (
		err error
	)

	p.SoloosEnv = soloosEnv

	err = p.SolonnClient.Init(p.SoloosEnv, solonnPeer.ID)
	if err != nil {
		return err
	}

	err = p.SolodnClient.Init(p.SoloosEnv)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloosEnv,
		&p.SolonnClient, &p.SolodnClient,
		p.NetBlockDriver.PrepareNetBlockMetaData,
		nil, nil, nil,
	)
	if err != nil {
		return err
	}

	err = p.MemBlockDriver.Init(p.SoloosEnv, memBlockDriverOptions)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.SoloosEnv, &p.NetBlockDriver, &p.MemBlockDriver,
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
