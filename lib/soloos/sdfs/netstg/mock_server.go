package netstg

import (
	"soloos/sdfs/protocol"
	"soloos/snet/srpc"
	snettypes "soloos/snet/types"
	"soloos/util"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	MockServerAddr = "127.0.0.1:10020"
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

	p.srpcServer.RegisterService(
		"/NetBlock/PWrite",
		func(requestID uint64, requestContentLen, parameterLen uint32, conn *snettypes.Connection) error {
			var blockData = make([]byte, requestContentLen)
			util.AssertErrIsNil(conn.ReadAll(blockData))
			var o protocol.UploadJob
			o.Init(blockData[:parameterLen], flatbuffers.GetUOffsetT(blockData[:parameterLen]))
			var backends = make([]protocol.UploadJobBackend, o.TransferBackendsLength())
			for i := 0; i < len(backends); i++ {
				o.TransferBackends(&backends[i], i)
			}
			util.AssertErrIsNil(conn.SimpleResponse(requestID, nil))
			return nil
		})
	return nil
}

func (p *MockServer) Serve() error {
	return p.srpcServer.Serve()
}

func (p *MockServer) Close() error {
	return p.srpcServer.Close()
}
