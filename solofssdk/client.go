package solofssdk

import (
	"soloos/common/fsapi"
	"soloos/common/log"
	"soloos/common/solodbapi"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/solomqapi"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

type Client struct {
	*soloosbase.SoloosEnv

	memStg      *memstg.MemStg
	metaPosixFs metastg.PosixFs
	memPosixFs  memstg.PosixFs

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

	err = p.metaPosixFs.Init(p.SoloosEnv,
		dbConn,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
	)
	if err != nil {
		log.Warn("Solofs metaPosixFs Init error", err)
		return err
	}

	err = p.memPosixFs.Init(p.SoloosEnv,
		nameSpaceID,
		p.memStg,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.metaPosixFs.FsINodeDriver.AllocFsINodeID,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
		p.metaPosixFs.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.metaPosixFs.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.metaPosixFs.FsINodeDriver.UpdateFsINodeInDB,
		p.metaPosixFs.FsINodeDriver.InsertFsINodeInDB,
		p.metaPosixFs.FsINodeDriver.FetchFsINodeByIDFromDB,
		p.metaPosixFs.FsINodeDriver.FetchFsINodeByNameFromDB,
		p.metaPosixFs.FIXAttrDriver.DeleteFIXAttrInDB,
		p.metaPosixFs.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.metaPosixFs.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		log.Warn("Solofs metaPosixFs Init error", err)
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.metaPosixFs.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) GetPosixFs() fsapi.PosixFs {
	return &p.memPosixFs
}
