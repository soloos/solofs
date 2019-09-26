package solofssdk

import (
	"soloos/common/log"
	"soloos/common/solodbapi"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
)

type ClientDriver struct {
	*soloosbase.SoloosEnv

	memStg memstg.MemStg
	dbConn solodbapi.Connection
}

var _ = solofsapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	solonnSRPCPeerID snettypes.PeerID,
	dbDriver string, dsn string,
) error {
	var err error
	p.SoloosEnv = soloosEnv

	err = p.initMemStg(solonnSRPCPeerID)
	if err != nil {
		log.Warn("solofs ClientDriver initMemStg error", err)
		return err
	}

	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		log.Warn("solofs ClientDriver dbConn init error", err)
		return err
	}

	return nil
}

func (p *ClientDriver) initMemStg(solonnSRPCPeerID snettypes.PeerID) error {
	var (
		err error
	)

	var solonnPeer snettypes.Peer
	solonnPeer, err = p.SoloosEnv.SNetDriver.GetPeer(solonnSRPCPeerID)
	if err != nil {
		log.Warn("solofs SNetDriver get solonnPeer error", err, solonnSRPCPeerID.Str())
		return err
	}

	err = p.memStg.Init(p.SoloosEnv, solonnPeer, memstg.MemBlockDriverOptions{})
	if err != nil {
		log.Warn("solofs memstg Init error", err)
		return err
	}

	return nil
}

func (p *ClientDriver) InitClient(itClient solofsapi.Client,
	nameSpaceID solofsapitypes.NameSpaceID,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	defaultMemBlocksLimit int32,
) error {
	var err error
	client := itClient.(*Client)

	err = p.memStg.MemBlockDriver.PrepareMemBlockTable(memstg.MemBlockTableOptions{
		ObjectSize:   defaultMemBlockCap,
		ObjectsLimit: defaultMemBlocksLimit,
	})
	if err != nil {
		log.Warn("Solofs ClientDriver PrepareMemBlockTabl error", err)
		return err
	}

	err = client.Init(p.SoloosEnv,
		nameSpaceID,
		&p.memStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
	if err != nil {
		log.Warn("Solofs ClientDriver InitClient error", err)
		return err
	}

	return nil
}
