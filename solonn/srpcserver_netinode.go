package solonn

import (
	"soloos/common/snet"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
)

func (p *SrpcServer) doNetINodeGet(isMustGet bool,
	reqCtx *snet.SNetReqContext,
	req *solofsprotocol.NetINodeInfoReq,
) (solofsprotocol.NetINodeInfoResp, error) {
	var (
		uNetINode  solofstypes.NetINodeUintptr
		netINodeID solofstypes.NetINodeID
		resp       solofsprotocol.NetINodeInfoResp
		err        error
	)

	netINodeID = req.NetINodeID
	if isMustGet {
		uNetINode, err = p.solonn.netINodeDriver.MustGetNetINode(netINodeID,
			req.Size, int(req.NetBlockCap), int(req.MemBlockCap))
	} else {
		uNetINode, err = p.solonn.netINodeDriver.GetNetINode(netINodeID)
	}
	defer p.solonn.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		if err.Error() == solofstypes.ErrObjectNotExists.Error() {
			return resp, nil
		}
		return resp, err
	}

	// response
	resp.Size = uNetINode.Ptr().Size
	resp.NetBlockCap = int32(uNetINode.Ptr().NetBlockCap)
	resp.MemBlockCap = int32(uNetINode.Ptr().MemBlockCap)

	return resp, nil
}

func (p *SrpcServer) NetINodeGet(reqCtx *snet.SNetReqContext,
	req solofsprotocol.NetINodeInfoReq,
) (solofsprotocol.NetINodeInfoResp, error) {
	return p.doNetINodeGet(false, reqCtx, &req)
}

func (p *SrpcServer) NetINodeMustGet(reqCtx *snet.SNetReqContext,
	req solofsprotocol.NetINodeInfoReq,
) (solofsprotocol.NetINodeInfoResp, error) {
	return p.doNetINodeGet(true, reqCtx, &req)
}

func (p *SrpcServer) NetINodeCommitSizeInDB(reqCtx *snet.SNetReqContext,
	req solofsprotocol.NetINodeCommitSizeInDBReq,
) error {
	var (
		uNetINode  solofstypes.NetINodeUintptr
		netINodeID solofstypes.NetINodeID
		err        error
	)

	netINodeID = req.NetINodeID
	uNetINode, err = p.solonn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solonn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	err = p.solonn.metaStg.NetINodeDriver.NetINodeTruncate(uNetINode, req.Size)
	if err != nil {
		return err
	}

	return nil
}
