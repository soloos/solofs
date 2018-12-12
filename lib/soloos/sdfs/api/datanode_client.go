package api

import (
	"soloos/snet"
)

type DataNodeClient struct {
	snetClientDriver *snet.ClientDriver
}

func (p *DataNodeClient) Init(snetClientDriver *snet.ClientDriver) error {
	p.snetClientDriver = snetClientDriver
	return nil
}
