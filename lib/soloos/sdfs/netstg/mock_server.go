package netstg

import (
	"soloos/log"
	"soloos/snet/srpc"
	snettypes "soloos/snet/types"
	"soloos/util"
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
		func(requestID uint64, requestContentLen uint32, conn *snettypes.Connection) error {
			log.Info("access")
			var blockData = make([]byte, requestContentLen)
			util.AssertErrIsNil(conn.ReadAll(blockData))
			util.AssertErrIsNil(conn.Response(requestID, nil))
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
