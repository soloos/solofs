package solodn

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"time"
)

func (p *Solodn) SetHeartBeatServers(heartBeatServerOptionsArr []snet.HeartBeatServerOptions) error {
	p.heartBeatServerOptionsArr = heartBeatServerOptionsArr
	return nil
}

func (p *Solodn) doHeartBeat(options snet.HeartBeatServerOptions) {
	var (
		heartBeat solofstypes.SolodnHeartBeat
		webret    iron.ResponseJSON
		peer      snet.Peer
		urlPath   string
		err       error
	)

	heartBeat.SrpcPeerID = p.srpcPeer.PeerID().Str()
	heartBeat.WebPeerID = p.webPeer.PeerID().Str()

	for {
		peer, err = p.SoloosEnv.SNetDriver.GetPeer(options.PeerID)
		urlPath = peer.AddressStr() + "/Api/Solofs/Solodn/HeartBeat"
		if err != nil {
			log.Error("Solodn HeartBeat post json error, urlPath:", urlPath, ", err:", err)
			goto HEARTBEAT_DONE
		}

		err = iron.PostJSON(urlPath, heartBeat, &webret)
		if err != nil {
			log.Error("Solodn HeartBeat post json(decode) error, urlPath:", urlPath, ", err:", err)
			goto HEARTBEAT_DONE
		}
		log.Info("Solodn heartbeat, urlPath:", urlPath, ", message:", webret)

	HEARTBEAT_DONE:
		time.Sleep(time.Duration(options.DurationMS) * time.Millisecond)
	}
}

func (p *Solodn) StartHeartBeat() error {
	for _, options := range p.heartBeatServerOptionsArr {
		go p.doHeartBeat(options)
	}
	return nil
}
