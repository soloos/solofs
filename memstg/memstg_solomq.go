package memstg

import "soloos/common/solomqapi"

func (p *MemStg) SetSolomqClient(solomqClient solomqapi.Client) error {
	p.solomqClient = solomqClient
	p.SolodnClient.SetUploadMemBlockWithSolomq(p.solomqClient.UploadMemBlockWithSolomq)
	p.NetBlockDriver.SetSolomqClient(solomqClient)
	return nil
}
