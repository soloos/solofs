package memstg

import "soloos/common/solomqapi"

func (p *MemStg) SetSOLOMQClient(solomqClient solomqapi.Client) error {
	p.solomqClient = solomqClient
	p.SolodnClient.SetUploadMemBlockWithSOLOMQ(p.solomqClient.UploadMemBlockWithSOLOMQ)
	p.NetBlockDriver.SetSOLOMQClient(solomqClient)
	return nil
}
