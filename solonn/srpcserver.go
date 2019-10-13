package solonn

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

type SrpcServer struct {
	solonn               *Solonn
	srpcServerListenAddr string
	srpcServerServeAddr  string
	srpcServer           snet.SrpcServer
}

var _ = iron.IServer(&SrpcServer{})

func (p *SrpcServer) Init(solonn *Solonn,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
) error {
	var err error

	p.solonn = solonn
	p.srpcServerListenAddr = srpcServerListenAddr
	p.srpcServerServeAddr = srpcServerServeAddr
	err = p.srpcServer.Init(solofstypes.DefaultSolofsRPCNetwork, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/Solodn/Register", p.SolodnRegister)

	p.srpcServer.RegisterService("/NetINode/Get", p.NetINodeGet)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetINode/CommitSizeInDB", p.NetINodeCommitSizeInDB)

	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)

	p.srpcServer.RegisterService("/FsINode/AllocFsINodeIno", p.solonn.metaStg.AllocFsINodeIno)
	p.srpcServer.RegisterService("/FsINode/DeleteFsINodeByIDInDB", p.solonn.metaStg.DeleteFsINodeByIDInDB)
	p.srpcServer.RegisterService("/FsINode/UpdateFsINodeInDB", p.solonn.metaStg.UpdateFsINodeInDB)
	p.srpcServer.RegisterService("/FsINode/InsertFsINodeInDB", p.solonn.metaStg.InsertFsINodeInDB)
	p.srpcServer.RegisterService("/FsINode/FetchFsINodeByIDFromDB", p.solonn.metaStg.FetchFsINodeByIDFromDB)
	p.srpcServer.RegisterService("/FsINode/FetchFsINodeByNameFromDB", p.solonn.metaStg.FetchFsINodeByNameFromDB)
	p.srpcServer.RegisterService("/FsINode/ListFsINodeByParentIDSelectCountFromDB",
		p.solonn.metaStg.ListFsINodeByParentIDSelectCountFromDB)
	p.srpcServer.RegisterService("/FsINode/ListFsINodeByParentIDSelectDataFromDB",
		p.solonn.metaStg.ListFsINodeByParentIDSelectDataFromDB)

	p.srpcServer.RegisterService("/FIXAttr/DeleteFIXAttrInDB", p.solonn.metaStg.DeleteFIXAttrInDB)
	p.srpcServer.RegisterService("/FIXAttr/ReplaceFIXAttrInDB", p.solonn.metaStg.ReplaceFIXAttrInDB)
	p.srpcServer.RegisterService("/FIXAttr/GetFIXAttrByInoFromDB", p.solonn.metaStg.GetFIXAttrByInoFromDB)

	return nil
}

func (p *SrpcServer) ServerName() string {
	return "Soloos.Solofs.Solonn.SrpcServer"
}

func (p *SrpcServer) Serve() error {
	log.Info("solonn srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *SrpcServer) Close() error {
	log.Info("solonn srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
