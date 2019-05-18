package libsdfs

import "soloos/common/swalapi"

func (p *Client) SetSWALClient(itSWALClient interface{}) error {
	var err error
	p.swalClient = itSWALClient.(swalapi.Client)

	err = p.memStg.RegisterSWALClient(p.swalClient)
	if err != nil {
		return err
	}

	return nil
}
