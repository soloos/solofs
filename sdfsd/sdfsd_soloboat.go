package sdfsd

func (p *SDFSD) initSoloBoat() error {
	return p.soloboatClient.Init(p.options.SoloBoatServeAddr)
}
