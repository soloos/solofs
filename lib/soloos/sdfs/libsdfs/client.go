package libsdfs

import (
	"soloos/dbcli"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/util/offheap"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver

	MemStg        *memstg.MemStg
	DirTreeDriver metastg.DirTreeDriver
	FdTable       FdTable
}

func (p *Client) Init(memStg *memstg.MemStg, dbConn *dbcli.Connection) error {
	var err error
	p.offheapDriver = &offheap.DefaultOffheapDriver

	p.MemStg = memStg

	err = p.DirTreeDriver.Init(p.offheapDriver,
		dbConn,
		p.MemStg.GetNetINodeWithReadAcquire,
		p.MemStg.MustGetNetINodeWithReadAcquire,
	)
	if err != nil {
		return err
	}

	err = p.FdTable.Init()
	if err != nil {
		return err
	}

	return nil
}
