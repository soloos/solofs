package main

import (
	"soloos/common/log"
	"time"
)

func (p *Env) initSoloBoat() error {
	p.soloboatCronJobDuration = time.Second * 3
	return p.soloboatClient.Init(p.options.SoloBoatServeAddr)
}

func (p *Env) doHeartBeat() error {
	return p.soloboatClient.HeartBeat(p.srpcPeerID)
}

func (p *Env) cronSoloBoatJob() error {
	go func() {
		for {
			var err = p.doHeartBeat()
			if err != nil {
				log.Warn("cronSoloBoatJob failed, err:", err)
			} else {
				log.Warn("cronSoloBoatJob HeartBeat, peerID:", p.srpcPeerID.Str())
			}
			time.Sleep(p.soloboatCronJobDuration)
		}
	}()
	return nil
}
