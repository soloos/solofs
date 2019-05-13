package netstg

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/common/snet/srpc"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
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
	srpcServer    srpc.Server
	dataNodePeers []snettypes.PeerUintptr
}

func (p *MockServer) SetDataNodePeers(dataNodePeers []snettypes.PeerUintptr) {
	p.dataNodePeers = dataNodePeers
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
	p.dataNodePeers = make([]snettypes.PeerUintptr, 3)
	for i := 0; i < len(p.dataNodePeers); i++ {
		p.dataNodePeers[i] = p.SNetDriver.AllocPeer(p.addr, sdfsapitypes.DefaultSDFSRPCProtocol)
	}

	return nil
}

func (p *MockServer) DataNodeRegister(serviceReq *snettypes.NetQuery) error {

	var param = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(param))

	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodeMustGet(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req protocol.NetINodeInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetNetINodeInfoResponse(&protocolBuilder, req.Size(), req.NetBlockCap(), req.MemBlockCap())
	util.AssertErrIsNil(
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodePWrite(serviceReq *snettypes.NetQuery) error {

	var reqBody = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqBody))

	var req protocol.NetINodePWriteRequest
	req.Init(reqBody[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqBody[:serviceReq.ParamSize]))
	var backends = make([]protocol.SNetPeer, req.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		req.TransferBackends(&backends[i], i)
	}

	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.SimpleResponse(serviceReq.ReqID, respBody))

	return nil
}

func (p *MockServer) NetINodePRead(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	var req protocol.NetINodePReadRequest
	req.Init(reqData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqData[:serviceReq.ParamSize]))

	var protocolBuilder flatbuffers.Builder
	api.SetNetINodePReadResponse(&protocolBuilder, req.Length())

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(serviceReq.Response(serviceReq.ReqID, respBody, make([]byte, req.Length())))
	return nil
}

func (p *MockServer) NetINodeCommitSizeInDB(serviceReq *snettypes.NetQuery) error {

	var reqData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(reqData))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

	return nil
}

func (p *MockServer) NetBlockPrepareMetaData(serviceReq *snettypes.NetQuery) error {

	var blockData = make([]byte, serviceReq.BodySize)
	util.AssertErrIsNil(serviceReq.ReadAll(blockData))

	// request
	var req protocol.NetINodeNetBlockInfoRequest
	req.Init(blockData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(blockData[:serviceReq.ParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetNetINodeNetBlockInfoResponse(&protocolBuilder, p.dataNodePeers[:], req.Cap(), req.Cap())
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
