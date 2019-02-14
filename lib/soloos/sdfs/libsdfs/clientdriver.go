package libsdfs

import (
	"soloos/common/sdbapi"
	"soloos/sdfs/memstg"
	"soloos/common/sdfsapi"
	"soloos/sdbone/offheap"
)

type ClientDriver struct {
	memStg memstg.MemStg
	dbConn sdbapi.Connection
}

var _ = sdfsapi.ClientDriver(&ClientDriver{})

func (p *ClientDriver) Init(nameNodeSRPCServerAddr string,
	defaultMemBlockChunkSize int, defaultMemBlockChunksLimit int32,
	dbDriver string, dsn string,
) error {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		err           error
	)

	err = initMemStg(&p.memStg, offheapDriver, nameNodeSRPCServerAddr, defaultMemBlockChunkSize, defaultMemBlockChunksLimit)
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
	defaultMemBlockChunkSize int, defaultMemBlockChunksLimit int32,
) error {
	var (
		memBlockDriverOptions = memstg.MemBlockDriverOptions{
			MemBlockPoolOptionsList: []memstg.MemBlockPoolOptions{
				memstg.MemBlockPoolOptions{
					defaultMemBlockChunkSize,
					defaultMemBlockChunksLimit,
				},
			},
		}
		err error
	)

	err = memStg.Init(offheapDriver, nameNodeSRPCServerAddr, memBlockDriverOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientDriver) InitClient(itClient sdfsapi.Client,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	client := itClient.(*Client)
	return client.Init(&p.memStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
}
