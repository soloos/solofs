package libsdfs

import (
	"soloos/common/sdbapi"
	"soloos/common/sdfsapi"
	"soloos/sdbone/offheap"
	"soloos/sdfs/memstg"
)

type ClientDriver struct {
	memStg memstg.MemStg
	dbConn sdbapi.Connection
}

var _ = sdfsapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(nameNodeSRPCServerAddr string,
	dbDriver string, dsn string,
) error {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		err           error
	)

	err = initMemStg(&p.memStg, offheapDriver, nameNodeSRPCServerAddr)
	if err != nil {
		return err
	}

	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		return err
	}

	return nil
}

func initMemStg(memStg *memstg.MemStg,
	offheapDriver *offheap.OffheapDriver,
	nameNodeSRPCServerAddr string,
) error {
	var (
		err error
	)

	err = memStg.Init(offheapDriver, nameNodeSRPCServerAddr, memstg.MemBlockDriverOptions{})
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

	err = client.Init(&p.memStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
	if err != nil {
		return err
	}

	return nil
}
