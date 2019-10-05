package solonn

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *SrpcServer) NetBlockPrepareMetaData(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodeNetBlockInfoReq,
) (solofsprotocol.NetINodeNetBlockInfoResp, error) {
	var (
		resp       solofsprotocol.NetINodeNetBlockInfoResp
		uNetINode  solofsapitypes.NetINodeUintptr
		netINodeID solofsapitypes.NetINodeID
		uNetBlock  solofsapitypes.NetBlockUintptr
		err        error
	)

	// request
	netINodeID = req.NetINodeID
	uNetINode, err = p.solonn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solonn.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		return resp, err
	}

	// response
	uNetBlock, err = p.solonn.netBlockDriver.MustGetNetBlock(uNetINode, req.NetBlockIndex)
	defer p.solonn.netBlockDriver.ReleaseNetBlock(uNetBlock)
	if err != nil {
		return resp, err
	}

	resp.Len = req.Cap
	resp.Cap = req.Cap
	resp.Backends = resp.Backends[:0]
	var peerIDs = uNetBlock.Ptr().StorDataBackends.Slice()
	for i, _ := range peerIDs {
		resp.Backends = append(resp.Backends, peerIDs[i].Str())
	}

	return resp, nil
}
