package memstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
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
	dataNodePeers []snettypes.Peer
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

	p.srpcServer.RegisterService("/DataNode/Register", p.DataNodeRegister)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)
	p.srpcServer.RegisterService("/NetINode/CommitSizeInDB", p.NetINodeCommitSizeInDB)
	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)
	p.dataNodePeers = make([]snettypes.Peer, 3)
	for i := 0; i < len(p.dataNodePeers); i++ {
		p.SNetDriver.InitPeerID((*snettypes.PeerID)(&p.dataNodePeers[i].ID))
		p.dataNodePeers[i].SetAddress(p.addr)
		p.dataNodePeers[i].ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
		util.AssertErrIsNil(p.SNetDriver.RegisterPeer(p.dataNodePeers[i]))
	}

	return nil
}

func (p *MockServer) DataNodeRegister(serviceReq *snettypes.NetQuery) error {

	var param = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(param))

	var protocolBuilder flatbuffers.Builder
	sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodeMustGet(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req sdfsprotocol.NetINodeInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	sdfsapi.SetNetINodeInfoResponse(&protocolBuilder, req.Size(), req.NetBlockCap(), req.MemBlockCap())
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodePWrite(serviceReq *snettypes.NetQuery) error {

	var reqBody = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqBody))

	var req sdfsprotocol.NetINodePWriteRequest
	req.Init(reqBody[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqBody[:serviceReq.ParamSize]))
	var backends = make([]string, req.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		backends[i] = string(req.TransferBackends(i))
	}

	var protocolBuilder flatbuffers.Builder
	sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.SimpleResponse(serviceReq.ReqID, respBody))

	return nil
}

func (p *MockServer) NetINodePRead(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	var req sdfsprotocol.NetINodePReadRequest
	req.Init(reqData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqData[:serviceReq.ParamSize]))

	var protocolBuilder flatbuffers.Builder
	sdfsapi.SetNetINodePReadResponse(&protocolBuilder, req.Length())

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.Response(serviceReq.ReqID, respBody, make([]byte, req.Length())))
	return nil
}

func (p *MockServer) NetINodeCommitSizeInDB(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	// response
	var protocolBuilder flatbuffers.Builder
	sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

	return nil
}

func (p *MockServer) NetBlockPrepareMetaData(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req sdfsprotocol.NetINodeNetBlockInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	var peerIDs = make([]snettypes.PeerID, len(p.dataNodePeers))
	for index, _ := range peerIDs {
		peerIDs[index] = p.dataNodePeers[index].PeerID()
	}
	sdfsapi.SetNetINodeNetBlockInfoResponse(&protocolBuilder, peerIDs, req.Cap(), req.Cap())
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
