package sdfssdk

import "soloos/common/swalapi"

func (p *Client) SetSWALClient(itSWALClient interface{}) error {
	var err error
	p.swalClient = itSWALClient.(swalapi.Client)

	err = p.memStg.SetSWALClient(p.swalClient)
	if err != nil {
		return err
	}

	return nil
}
