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

	memStg     *memstg.MemStg
	memPosixFs memstg.PosixFs

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

	err = p.memPosixFs.Init(p.SoloosEnv,
		nsID,
		p.memStg,
		defaultNetBlockCap,
		defaultMemBlockCap,
		p.memStg.GetNetINode,
		p.memStg.MustGetNetINode,
		p.memStg.ReleaseNetINode,
		p.memPosixFs.FsINodeDriver.AllocFsINodeID,
		p.memPosixFs.FsINodeDriver.DeleteFsINodeByIDInDB,
		p.memPosixFs.FsINodeDriver.ListFsINodeByParentIDFromDB,
		p.memPosixFs.FsINodeDriver.UpdateFsINodeInDB,
		p.memPosixFs.FsINodeDriver.InsertFsINodeInDB,
		p.memPosixFs.FsINodeDriver.FetchFsINodeByIDFromDB,
		p.memPosixFs.FsINodeDriver.FetchFsINodeByNameFromDB,
		p.memPosixFs.FIXAttrDriver.DeleteFIXAttrInDB,
		p.memPosixFs.FIXAttrDriver.ReplaceFIXAttrInDB,
		p.memPosixFs.FIXAttrDriver.GetFIXAttrByInoFromDB,
	)
	if err != nil {
		log.Warn("Solofs memPosixFs Init error", err)
		return err
	}

	return nil
}

func (p *Client) Close() error {
	var err error
	err = p.memPosixFs.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) GetPosixFs() fsapi.PosixFs {
	return &p.memPosixFs
}
