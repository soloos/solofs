package solofssdk

import "soloos/common/solomqapi"

func (p *Client) SetSOLOMQClient(itSOLOMQClient interface{}) error {
	var err error
	p.solomqClient = itSOLOMQClient.(solomqapi.Client)

	err = p.memStg.SetSOLOMQClient(p.solomqClient)
	if err != nil {
		return err
	}

	return nil
}
