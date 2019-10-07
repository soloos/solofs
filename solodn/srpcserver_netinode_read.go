package solodn

import (
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
)

func (p *SrpcServer) NetINodePRead(reqCtx *snet.SNetReqContext,
	req solofsprotocol.NetINodePWriteReq,
) error {
	var (
		uNetBlock          solofstypes.NetBlockUintptr
		netINodeID         solofstypes.NetINodeID
		uNetINode          solofstypes.NetINodeUintptr
		resp               solofsprotocol.NetINodePReadResp
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
		readDataSize       int
		err                error
	)
	netINodeID = req.NetINodeID
	uNetINode, err = p.solodn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solodn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	// calculate file data size
	if req.Offset+uint64(req.Length) > uNetINode.Ptr().Size {
		readDataSize = int(uNetINode.Ptr().Size - req.Offset)
	} else {
		readDataSize = int(req.Length)
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(req.Offset / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((req.Offset + uint64(readDataSize)) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.solodn.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.solodn.netBlockDriver.ReleaseNetBlock(uNetBlock)
		if err != nil {
			return err
		}

		if uNetBlock.Ptr().IsLocalDataBackendInited.Load() == solodbtypes.MetaDataStateUninited {
			p.solodn.PrepareNetBlockLocalDataBackend(uNetBlock)
		}
	}

	// request file data
	reqCtx.SetResponseInService()
	resp.Length = int32(readDataSize)
	err = reqCtx.ResponseWithOffheap(reqCtx.ReqID,
		&resp, readDataSize,
	)
	if err != nil {
		log.Debug("NetINodePRead SimpleResponse error,err:", err)
		return nil
	}

	_, err = p.solodn.netINodeDriver.PReadWithNetQuery(uNetINode, &reqCtx.NetQuery,
		int(readDataSize), req.Offset)
	if err != nil {
		log.Debug("NetINodePRead SimpleResponse error,err:", err)
		return nil
	}

	return nil
}
