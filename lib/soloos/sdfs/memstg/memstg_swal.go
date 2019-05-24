package memstg

import "soloos/common/swalapi"

func (p *MemStg) SetSWALClient(swalClient swalapi.Client) error {
	p.swalClient = swalClient
	p.DataNodeClient.SetUploadMemBlockWithSWAL(p.swalClient.UploadMemBlockWithSWAL)
	p.NetBlockDriver.SetSWALClient(swalClient)
	return nil
}
