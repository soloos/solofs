package sdfssdk

import (
	"soloos/common/log"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/memstg"
)

type ClientDriver struct {
	*soloosbase.SoloOSEnv

	memStg memstg.MemStg
	dbConn sdbapi.Connection
}

var _ = sdfsapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodePeerID snettypes.PeerID,
	dbDriver string, dsn string,
) error {
	var err error
	p.SoloOSEnv = soloOSEnv

	err = p.initMemStg(nameNodePeerID)
	if err != nil {
		log.Debug("sdfs ClientDriver initMemStg error", err)
		return err
	}

	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		log.Debug("sdfs ClientDriver dbConn init error", err)
		return err
	}

	return nil
}

func (p *ClientDriver) initMemStg(nameNodePeerID snettypes.PeerID) error {
	var (
		err error
	)

	var nameNodePeer snettypes.Peer
	nameNodePeer, err = p.SoloOSEnv.SNetDriver.GetPeer(nameNodePeerID)
	if err != nil {
		log.Debug("sdfs SNetDriver get nameNodePeer error", err, nameNodePeerID.Str())
		return err
	}

	err = p.memStg.Init(p.SoloOSEnv, nameNodePeer, memstg.MemBlockDriverOptions{})
	if err != nil {
		log.Debug("sdfs memstg Init error", err)
		return err
	}

	return nil
}

func (p *ClientDriver) InitClient(itClient sdfsapi.Client,
	nameSpaceID sdfsapitypes.NameSpaceID,
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
		log.Debug("SDFS ClientDriver PrepareMemBlockTabl error", err)
		return err
	}

	err = client.Init(p.SoloOSEnv,
		nameSpaceID,
		&p.memStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
	if err != nil {
		log.Debug("SDFS ClientDriver InitClient error", err)
		return err
	}

	return nil
}
