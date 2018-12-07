package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/snet/srpc"
	snettypes "soloos/snet/types"
	"soloos/util"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	DefaultMockServerAddr = "127.0.0.1:10020"
)

type MockServer struct {
	snetDriver    *snet.SNetDriver
	network       string
	addr          string
	srpcServer    srpc.Server
	dataNodePeers [3]snettypes.PeerUintptr
}

func (p *MockServer) Init(snetDriver *snet.SNetDriver, network string, addr string) error {
	var err error
	p.snetDriver = snetDriver
	p.network = network
	p.addr = addr
	err = p.srpcServer.Init(p.network, p.addr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetBlock/PWrite", p.NetBlockPWrite)
	p.srpcServer.RegisterService("/NetBlock/PRead", p.NetBlockPRead)
	p.srpcServer.RegisterService("/NetBlock/MustGet", p.NetBlockMustGet)
	for i := 0; i < len(p.dataNodePeers); i++ {
		p.dataNodePeers[i] = p.snetDriver.MustGetPeer(nil, p.addr, types.DefaultSDFSRPCProtocol)
	}

	return nil
}

func (p *MockServer) NetBlockPWrite(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var reqBody = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(reqBody))

	var req protocol.NetBlockPWriteRequest
	req.Init(reqBody[:reqParamSize], flatbuffers.GetUOffsetT(reqBody[:reqParamSize]))
	var backends = make([]protocol.NetBlockBackend, req.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		req.TransferBackends(&backends[i], i)
	}

	var protocolBuilder flatbuffers.Builder
	protocol.CommonResponseStart(&protocolBuilder)
	protocol.CommonResponseAddCode(&protocolBuilder, snettypes.CODE_OK)
	protocolBuilder.Finish(protocol.CommonResponseEnd(&protocolBuilder))
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(conn.SimpleResponse(reqID, respBody))

	return nil
}

func (p *MockServer) NetBlockPRead(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var reqData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(reqData))

	var req protocol.NetBlockPReadRequest
	req.Init(reqData[:reqParamSize], flatbuffers.GetUOffsetT(reqData[:reqParamSize]))

	var protocolBuilder flatbuffers.Builder
	api.SetNetBlockPReadResponse(snettypes.CODE_OK, req.Length(), &protocolBuilder)

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(conn.Response(reqID, respBody, make([]byte, req.Length())))
	return nil
}

func (p *MockServer) NetBlockMustGet(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var blockData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(blockData))

	// request
	var req protocol.INodeNetBlockInfoRequest
	req.Init(blockData[:reqParamSize], flatbuffers.GetUOffsetT(blockData[:reqParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetINodeNetBlockInfoResp(p.dataNodePeers[:], req.Cap(), req.Cap(), &protocolBuilder)
	util.AssertErrIsNil(conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

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
