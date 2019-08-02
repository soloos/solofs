package datanode

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"time"
)

func (p *DataNode) SetHeartBeatServers(heartBeatServerOptionsArr []sdfsapitypes.HeartBeatServerOptions) error {
	p.heartBeatServerOptionsArr = heartBeatServerOptionsArr
	return nil
}

func (p *DataNode) doHeartBeat(options sdfsapitypes.HeartBeatServerOptions) {
	var (
		heartBeat sdfsapitypes.DataNodeHeartBeat
		webret    iron.ApiOutputResult
		peer      snettypes.Peer
		urlPath   string
		err       error
	)

	heartBeat.SRPCPeerID = p.srpcPeer.PeerID().Str()
	heartBeat.WebPeerID = p.webPeer.PeerID().Str()

	for {
		peer, err = p.SoloOSEnv.SNetDriver.GetPeer(options.PeerID)
		urlPath = peer.AddressStr() + "/Api/SDFS/DataNode/HeartBeat"
		if err != nil {
			log.Error("DataNode HeartBeat post json error, urlPath:", urlPath, ", err:", err)
			goto HEARTBEAT_DONE
		}

		err = iron.PostJSON(urlPath, heartBeat, &webret)
		if err != nil {
			log.Error("DataNode HeartBeat post json(decode) error, urlPath:", urlPath, ", err:", err)
			goto HEARTBEAT_DONE
		}
		log.Info("DataNode heartbeat, urlPath:", urlPath, ", message:", webret)

	HEARTBEAT_DONE:
		time.Sleep(time.Duration(options.DurationMS) * time.Millisecond)
	}
}

func (p *DataNode) StartHeartBeat() error {
	for _, options := range p.heartBeatServerOptionsArr {
		go p.doHeartBeat(options)
	}
	return nil
}
