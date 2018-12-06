package api

import (
	"soloos/snet"
)

type DataNodeClient struct {
	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver
}

func (p *DataNodeClient) Init(snetDriver *snet.SNetDriver, snetClientDriver *snet.ClientDriver) error {
	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	return nil
}
