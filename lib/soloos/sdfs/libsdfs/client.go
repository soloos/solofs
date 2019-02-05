package libsdfs

import (
	"soloos/dbcli"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/util/offheap"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver

	MemStg         *memstg.MemStg
	MetaDirTreeStg metastg.DirTreeStg
	MemDirTreeStg  memstg.DirTreeStg
}

func (p *Client) Init(memStg *memstg.MemStg,
	dbConn *dbcli.Connection,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	var err error
	p.offheapDriver = &offheap.DefaultOffheapDriver

	p.MemStg = memStg

	err = p.MetaDirTreeStg.Init(p.offheapDriver,
		dbConn,
		p.MemStg.GetNetINodeWithReadAcquire,
		p.MemStg.MustGetNetINodeWithReadAcquire,
	)
	if err != nil {
		return err
	}

	err = p.MemDirTreeStg.SdfsInit(
		p.MemStg,
		p.offheapDriver,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.MetaDirTreeStg.FsINodeDriver.AllocFsINodeID,
		p.MemStg.GetNetINodeWithReadAcquire,
		p.MemStg.MustGetNetINodeWithReadAcquire,
		p.MetaDirTreeStg.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.MetaDirTreeStg.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.MetaDirTreeStg.FsINodeDriver.UpdateFsINodeInDB,
		p.MetaDirTreeStg.FsINodeDriver.InsertFsINodeInDB,
		p.MetaDirTreeStg.FsINodeDriver.GetFsINodeByIDFromDB,
		p.MetaDirTreeStg.FsINodeDriver.GetFsINodeByNameFromDB,
		p.MetaDirTreeStg.FIXAttrDriver.DeleteFIXAttrInDB,
		p.MetaDirTreeStg.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.MetaDirTreeStg.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.MetaDirTreeStg.Close()
	if err != nil {
		return err
	}

	return nil
}
