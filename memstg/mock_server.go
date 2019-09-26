package memstg

import (
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	DefaultMockServerAddr = "127.0.0.1:10020"
)

type MockServer struct {
	*soloosbase.SoloOSEnv
	network       string
	addr          string
	srpcServer    snet.SRPCServer
	solodnPeers []snettypes.Peer
}

func (p *MockServer) Init(soloOSEnv *soloosbase.SoloOSEnv, network string, addr string) error {
	var err error
	p.SoloOSEnv = soloOSEnv
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
		p.solodnPeers[i].ServiceProtocol = solofsapitypes.DefaultSOLOFSRPCProtocol
		util.AssertErrIsNil(p.SNetDriver.RegisterPeer(p.solodnPeers[i]))
	}

	return nil
}

func (p *MockServer) SolodnRegister(serviceReq *snettypes.NetQuery) error {

	var param = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(param))

	var protocolBuilder flatbuffers.Builder
	solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodeMustGet(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req solofsprotocol.NetINodeInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	solofsapi.SetNetINodeInfoResponse(&protocolBuilder, req.Size(), req.NetBlockCap(), req.MemBlockCap())
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodePWrite(serviceReq *snettypes.NetQuery) error {

	var reqBody = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqBody))

	var req solofsprotocol.NetINodePWriteRequest
	req.Init(reqBody[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqBody[:serviceReq.ParamSize]))
	var backends = make([]string, req.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		backends[i] = string(req.TransferBackends(i))
	}

	var protocolBuilder flatbuffers.Builder
	solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.SimpleResponse(serviceReq.ReqID, respBody))

	return nil
}

func (p *MockServer) NetINodePRead(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	var req solofsprotocol.NetINodePReadRequest
	req.Init(reqData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqData[:serviceReq.ParamSize]))

	var protocolBuilder flatbuffers.Builder
	solofsapi.SetNetINodePReadResponse(&protocolBuilder, req.Length())

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.Response(serviceReq.ReqID, respBody, make([]byte, req.Length())))
	return nil
}

func (p *MockServer) NetINodeCommitSizeInDB(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	// response
	var protocolBuilder flatbuffers.Builder
	solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

	return nil
}

func (p *MockServer) NetBlockPrepareMetaData(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req solofsprotocol.NetINodeNetBlockInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	var peerIDs = make([]snettypes.PeerID, len(p.solodnPeers))
	for index, _ := range peerIDs {
		peerIDs[index] = p.solodnPeers[index].PeerID()
	}
	solofsapi.SetNetINodeNetBlockInfoResponse(&protocolBuilder, peerIDs, req.Cap(), req.Cap())
	util.AssertErrIsNil(serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
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

func MakeMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string, mockServer *MockServer) {
	util.AssertErrIsNil(mockServer.Init(soloOSEnv, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
}
