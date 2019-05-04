package libsdfs

import (
	"soloos/common/fsapi"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
)

type Client struct {
	*soloosbase.SoloOSEnv

	memStg         *memstg.MemStg
	metaDirTreeStg metastg.DirTreeStg
	memDirTreeStg  memstg.DirTreeStg
}

var _ = sdfsapi.Client(&Client{})

func (p *Client) Init(soloOSEnv *soloosbase.SoloOSEnv,
	memStg *memstg.MemStg,
	dbConn *sdbapi.Connection,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.memStg = memStg

	err = p.metaDirTreeStg.Init(p.SoloOSEnv,
		dbConn,
		p.memStg.GetNetINodeWithReadAcquire,
		p.memStg.MustGetNetINodeWithReadAcquire,
	)
	if err != nil {
		return err
	}

	err = p.memDirTreeStg.Init(p.SoloOSEnv,
		p.memStg,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.metaDirTreeStg.FsINodeDriver.AllocFsINodeID,
		p.memStg.GetNetINodeWithReadAcquire,
		p.memStg.MustGetNetINodeWithReadAcquire,
		p.metaDirTreeStg.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.metaDirTreeStg.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.metaDirTreeStg.FsINodeDriver.UpdateFsINodeInDB,
		p.metaDirTreeStg.FsINodeDriver.InsertFsINodeInDB,
		p.metaDirTreeStg.FsINodeDriver.GetFsINodeByIDFromDB,
		p.metaDirTreeStg.FsINodeDriver.GetFsINodeByNameFromDB,
		p.metaDirTreeStg.FIXAttrDriver.DeleteFIXAttrInDB,
		p.metaDirTreeStg.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.metaDirTreeStg.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.metaDirTreeStg.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) GetPosixFS() fsapi.PosixFS {
	return &p.memDirTreeStg
}
