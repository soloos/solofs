package memstg

import "soloos/common/swalapi"

func (p *MemStg) RegisterSWALClient(swalClient swalapi.Client) error {
	p.swalClient = swalClient
	p.DataNodeClient.SetUploadMemBlockWithSWAL(p.swalClient.UploadMemBlockWithSWAL)
	return nil
}
