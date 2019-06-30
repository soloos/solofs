package main

import (
	"soloos/common/log"
	"time"
)

func (p *Env) initSilicon() error {
	p.siliconCronJobDuration = time.Second * 3
	return p.siliconClient.Init(p.options.SiliconServeAddr)
}

func (p *Env) doHeartBeat() error {
	return p.siliconClient.HeartBeat(p.peerID)
}

func (p *Env) cronSiliconJob() error {
	go func() {
		for {
			var err = p.doHeartBeat()
			if err != nil {
				log.Warn("cronSiliconJob failed, err:", err)
			} else {
				log.Warn("cronSiliconJob HeartBeat, peerID:", p.peerID.Str())
			}
			time.Sleep(p.siliconCronJobDuration)
		}
	}()
	return nil
}
