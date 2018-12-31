package libsdfs

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/util/offheap"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver

	MemStg  memstg.MemStg
	MetaStg metastg.MetaStg
}

func (p *Client) Init(nameNodeSRPCServerAddr string,
	memBlockChunkSize int, memBlockChunksLimit int32,
	dbDriver, dsn string,
) error {
	var err error
	p.offheapDriver = &offheap.DefaultOffheapDriver

	var memBlockDriverOptions = memstg.MemBlockDriverOptions{
		MemBlockPoolOptionsList: []memstg.MemBlockPoolOptions{
			memstg.MemBlockPoolOptions{
				memBlockChunkSize,
				memBlockChunksLimit,
			},
		},
	}
	err = p.MemStg.Init(p.offheapDriver, nameNodeSRPCServerAddr, memBlockDriverOptions)
	if err != nil {
		return err
	}

	err = p.MetaStg.Init(p.offheapDriver, dbDriver, dsn,
		p.MemStg.NetINodeDriver.GetNetINode,
		p.MemStg.NetINodeDriver.MustGetNetINode,
	)
	if err != nil {
		return err
	}

	return nil
}
