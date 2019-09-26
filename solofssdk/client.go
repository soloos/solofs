package solofssdk

import (
	"soloos/common/fsapi"
	"soloos/common/log"
	"soloos/common/solodbapi"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/solomqapi"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

type Client struct {
	*soloosbase.SoloosEnv

	memStg      *memstg.MemStg
	metaPosixFS metastg.PosixFS
	memPosixFS  memstg.PosixFS

	solomqClient solomqapi.Client
}

var _ = solofsapi.Client(&Client{})

func (p *Client) Init(soloosEnv *soloosbase.SoloosEnv,
	nameSpaceID solofsapitypes.NameSpaceID,
	memStg *memstg.MemStg,
	dbConn *solodbapi.Connection,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.memStg = memStg

	err = p.metaPosixFS.Init(p.SoloosEnv,
		dbConn,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
	)
	if err != nil {
		log.Warn("Solofs metaPosixFS Init error", err)
		return err
	}

	err = p.memPosixFS.Init(p.SoloosEnv,
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
		log.Warn("Solofs metaPosixFS Init error", err)
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
