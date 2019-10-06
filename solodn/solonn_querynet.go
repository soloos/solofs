package solodn

import "soloos/common/solofsprotocol"

func (p *Solodn) RegisterInSolonn() error {
	var req = solofsprotocol.SNetPeer{
		PeerID:   p.srpcPeer.PeerIDStr(),
		Address:  p.srpcPeer.AddressStr(),
		Protocol: p.srpcPeer.ServiceProtocol.Str(),
	}

	return p.solonnClient.Dispatch("/Solodn/Register", nil, req)
}
