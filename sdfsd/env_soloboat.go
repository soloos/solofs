package main

func (p *Env) initSoloBoat() error {
	return p.soloboatClient.Init(p.options.SoloBoatServeAddr)
}
