package libsdfs

import (
	"soloos/dbcli"
	"soloos/sdfs/memstg"
	"soloos/util/offheap"
)

type ClientDriver struct {
	MemStg memstg.MemStg
	dbConn dbcli.Connection
}

func (p *ClientDriver) Init(nameNodeSRPCServerAddr string,
	memBlockChunkSize int, memBlockChunksLimit int32,
	dbDriver string, dsn string,
) error {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		err           error
	)

	err = initMemStg(&p.MemStg, offheapDriver, nameNodeSRPCServerAddr, memBlockChunkSize, memBlockChunksLimit)
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
	memBlockChunkSize int, memBlockChunksLimit int32,
) error {
	var (
		memBlockDriverOptions = memstg.MemBlockDriverOptions{
			MemBlockPoolOptionsList: []memstg.MemBlockPoolOptions{
				memstg.MemBlockPoolOptions{
					memBlockChunkSize,
					memBlockChunksLimit,
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

func (p *ClientDriver) InitClient(client *Client,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	return client.Init(&p.MemStg, &p.dbConn,
		defaultNetBlockCap,
		defaultMemBlockCap,
	)
}
