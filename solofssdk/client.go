package solofssdk

import (
	"soloos/common/fsapi"
	"soloos/common/log"
	"soloos/common/solodbapi"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"soloos/common/solomqapi"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
)

type Client struct {
	*soloosbase.SoloosEnv

	memStg *memstg.MemStg
	memstg.PosixFs

	solomqClient solomqapi.Client
}

var _ = solofsapi.Client(&Client{})

func (p *Client) Init(soloosEnv *soloosbase.SoloosEnv,
	nsID solofstypes.NameSpaceID,
	memStg *memstg.MemStg,
	dbConn *solodbapi.Connection,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.memStg = memStg

	err = p.PosixFs.Init(p.SoloosEnv,
		nsID,
		p.memStg,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
		p.PosixFs.FsINodeDriver.AllocFsINodeID,
		p.PosixFs.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.PosixFs.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.PosixFs.FsINodeDriver.UpdateFsINodeInDB,
		p.PosixFs.FsINodeDriver.InsertFsINodeInDB,
		p.PosixFs.FsINodeDriver.FetchFsINodeByIDFromDB,
		p.PosixFs.FsINodeDriver.FetchFsINodeByNameFromDB,
		p.PosixFs.FIXAttrDriver.DeleteFIXAttrInDB,
		p.PosixFs.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.PosixFs.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		log.Warn("Solofs MemPosixFs Init error", err)
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.PosixFs.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) GetPosixFs() fsapi.PosixFs {
	return &p.PosixFs
}
