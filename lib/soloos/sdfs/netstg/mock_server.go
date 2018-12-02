package netstg

import (
	"soloos/sdfs/protocol"
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
	srpcServer srpc.Server
}

func (p *MockServer) Init(network string, addr string) error {
	var err error
	err = p.srpcServer.Init(network, addr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetBlock/PWrite", p.NetBlockPWriteRequest)
	return nil
}

func (p *MockServer) NetBlockPWriteRequest(requestID uint64,
	requestContentLen, parameterLen uint32,
	conn *snettypes.Connection) error {
	var blockData = make([]byte, requestContentLen)
	util.AssertErrIsNil(conn.ReadAll(blockData))
	var o protocol.NetBlockPWriteRequest
	o.Init(blockData[:parameterLen], flatbuffers.GetUOffsetT(blockData[:parameterLen]))
	var backends = make([]protocol.NetBlockBackend, o.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		o.TransferBackends(&backends[i], i)
	}
	util.AssertErrIsNil(conn.SimpleResponse(requestID, nil))
	return nil
}

func (p *MockServer) NetBlockPRead(requestID uint64,
	requestContentLen, parameterLen uint32,
	conn *snettypes.Connection) error {
	var blockData = make([]byte, requestContentLen)
	util.AssertErrIsNil(conn.ReadAll(blockData))
	var o protocol.NetBlockPWriteRequest
	o.Init(blockData[:parameterLen], flatbuffers.GetUOffsetT(blockData[:parameterLen]))
	var backends = make([]protocol.NetBlockBackend, o.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		o.TransferBackends(&backends[i], i)
	}
	util.AssertErrIsNil(conn.SimpleResponse(requestID, nil))
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
