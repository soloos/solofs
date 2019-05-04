package libsdfs

import (
	"soloos/common/sdbapi"
	"soloos/common/sdfsapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/memstg"
)

type ClientDriver struct {
	*soloosbase.SoloOSEnv

	memStg memstg.MemStg
	dbConn sdbapi.Connection
}

var _ = sdfsapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	dbDriver string, dsn string,
) error {
	var err error
	p.SoloOSEnv = soloOSEnv

	err = p.initMemStg(nameNodeSRPCServerAddr)
	if err != nil {
		return err
	}

	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) initMemStg(nameNodeSRPCServerAddr string) error {
	var (
		err error
	)

	err = p.memStg.Init(p.SoloOSEnv, nameNodeSRPCServerAddr, memstg.MemBlockDriverOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) InitClient(itClient sdfsapi.Client,
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
		return err
	}

	err = client.Init(p.SoloOSEnv, &p.memStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
	if err != nil {
		return err
	}

	return nil
}
