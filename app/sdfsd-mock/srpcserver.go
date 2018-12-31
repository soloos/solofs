package main

import (
	"fmt"
	"soloos/snet/types"
)

func (p *Env) initSRPCServer(srpcServerPort int) error {
	var (
		err error
	)
	err = p.SRPCServer.Init("tcp", fmt.Sprintf("0.0.0.0:%d", srpcServerPort))
	if err != nil {
		return err
	}

	p.SRPCServer.RegisterService("/NetBlock/Write", p.SRPCNetBlockWrite)

	return nil
}

func (p *Env) SRPCNetBlockWrite(requestID uint64, conn *types.Connection) error {
	return nil
}
