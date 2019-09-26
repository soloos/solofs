package solofssdk

import "soloos/common/solomqapi"

func (p *Client) SetSolomqClient(itSolomqClient interface{}) error {
	var err error
	p.solomqClient = itSolomqClient.(solomqapi.Client)

	err = p.memStg.SetSolomqClient(p.solomqClient)
	if err != nil {
		return err
	}

	return nil
}
