package memstg

import (
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"time"
)

const (
	DefaultMockServerAddr = "127.0.0.1:10020"
)

type MockServer struct {
	*soloosbase.SoloosEnv
	network     string
	addr        string
	srpcServer  snet.SrpcServer
	solodnPeers []snettypes.Peer
}

func (p *MockServer) Init(soloosEnv *soloosbase.SoloosEnv, network string, addr string) error {
	var err error
	p.SoloosEnv = soloosEnv
	p.network = network
	p.addr = addr
	err = p.srpcServer.Init(p.network, p.addr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/Solodn/Register", p.SolodnRegister)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)
	p.srpcServer.RegisterService("/NetINode/CommitSizeInDB", p.NetINodeCommitSizeInDB)
	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)
	p.solodnPeers = make([]snettypes.Peer, 3)
	for i := 0; i < len(p.solodnPeers); i++ {
		p.SNetDriver.InitPeerID((*snettypes.PeerID)(&p.solodnPeers[i].ID))
		p.solodnPeers[i].SetAddress(p.addr)
		p.solodnPeers[i].ServiceProtocol = solofsapitypes.DefaultSolofsRPCProtocol
		util.AssertErrIsNil(p.SNetDriver.RegisterPeer(p.solodnPeers[i]))
	}

	return nil
}

func (p *MockServer) SolodnRegister(reqCtx *snettypes.SNetReqContext) error {
	return nil
}

func (p *MockServer) NetINodeMustGet(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodeInfoReq,
) solofsprotocol.NetINodeInfoResp {
	util.AssertErrIsNil(reqCtx.SkipReadRemaining())

	var resp solofsprotocol.NetINodeInfoResp
	resp.Size = req.Size
	resp.NetBlockCap = req.NetBlockCap
	resp.MemBlockCap = req.MemBlockCap

	return resp
}

func (p *MockServer) NetINodePWrite(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodePWriteReq,
) error {
	util.AssertErrIsNil(reqCtx.SkipReadRemaining())
	return nil
}

func (p *MockServer) NetINodePRead(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodePReadReq,
) solofsprotocol.NetINodePReadResp {
	util.AssertErrIsNil(reqCtx.SkipReadRemaining())
	return solofsprotocol.NetINodePReadResp{Length: req.Length}
}

func (p *MockServer) NetINodeCommitSizeInDB(reqCtx *snettypes.SNetReqContext) error {
	util.AssertErrIsNil(reqCtx.SkipReadRemaining())
	return nil
}

func (p *MockServer) NetBlockPrepareMetaData(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodeNetBlockInfoReq,
) solofsprotocol.NetINodeNetBlockInfoResp {
	util.AssertErrIsNil(reqCtx.SkipReadRemaining())

	var resp solofsprotocol.NetINodeNetBlockInfoResp
	resp.Cap = req.Cap
	resp.Len = req.Cap
	resp.Backends = resp.Backends[:0]
	for _, peer := range p.solodnPeers {
		resp.Backends = append(resp.Backends, peer.PeerID().Str())
	}

	return resp
}

func (p *MockServer) Serve() error {
	return p.srpcServer.Serve()
}

func (p *MockServer) Close() error {
	var err error
	err = p.srpcServer.Close()
	time.Sleep(time.Second)
	return err
}

func MakeMockServerForTest(soloosEnv *soloosbase.SoloosEnv,
	mockServerAddr string, mockServer *MockServer) {
	util.AssertErrIsNil(mockServer.Init(soloosEnv, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
}
