package sdfscli

import (
	"soloos/common/fsapi"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/swalapi"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
)

type Client struct {
	*soloosbase.SoloOSEnv

	memStg      *memstg.MemStg
	metaPosixFS metastg.PosixFS
	memPosixFS  memstg.PosixFS

	swalClient swalapi.Client
}

var _ = sdfsapi.Client(&Client{})

func (p *Client) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameSpaceID sdfsapitypes.NameSpaceID,
	memStg *memstg.MemStg,
	dbConn *sdbapi.Connection,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.memStg = memStg

	err = p.metaPosixFS.Init(p.SoloOSEnv,
		dbConn,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
	)
	if err != nil {
		return err
	}

	err = p.memPosixFS.Init(p.SoloOSEnv,
		nameSpaceID,
		p.memStg,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.metaPosixFS.FsINodeDriver.AllocFsINodeID,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
		p.metaPosixFS.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.metaPosixFS.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.metaPosixFS.FsINodeDriver.UpdateFsINodeInDB,
		p.metaPosixFS.FsINodeDriver.InsertFsINodeInDB,
		p.metaPosixFS.FsINodeDriver.FetchFsINodeByIDFromDB,
		p.metaPosixFS.FsINodeDriver.FetchFsINodeByNameFromDB,
		p.metaPosixFS.FIXAttrDriver.DeleteFIXAttrInDB,
		p.metaPosixFS.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.metaPosixFS.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.metaPosixFS.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) GetPosixFS() fsapi.PosixFS {
	return &p.memPosixFS
}
