package libsdfs

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"sync"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver

	FileTableIndexs  sync.NoGCUintptrPool
	FileTableRWMutex sync.RWMutex
	FileTable        []types.FsINodeFileHandler

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

	p.FileTableIndexs.New = func() uintptr {
		var fdID uintptr
		p.FileTableRWMutex.Lock()
		fdID = uintptr(len(p.FileTable))
		p.FileTable = append(p.FileTable, types.FsINodeFileHandler{})
		p.FileTableRWMutex.Unlock()
		return fdID
	}

	return nil
}
