package solodn

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *SrpcServer) NetINodeSync(reqCtx *snet.SNetReqContext,
	req solofsprotocol.NetINodePWriteReq,
) error {
	var (
		netINodeID solofsapitypes.NetINodeID
		uNetINode  solofsapitypes.NetINodeUintptr
		err        error
	)
	netINodeID = req.NetINodeID
	uNetINode, err = p.solodn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solodn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	err = p.solodn.netINodeDriver.Sync(uNetINode)
	if err != nil {
		return err
	}

	return nil

}
