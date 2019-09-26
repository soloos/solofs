package solodn

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) NetINodePRead(serviceReq *snettypes.NetQuery) error {
	var (
		reqParamData = make([]byte, serviceReq.ParamSize)
		reqParam     solofsprotocol.NetINodePWriteRequest
		uNetBlock    solofsapitypes.NetBlockUintptr
		err          error
	)

	// request param
	err = serviceReq.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:serviceReq.ParamSize],
		flatbuffers.GetUOffsetT(reqParamData[:serviceReq.ParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder    flatbuffers.Builder
		netINodeID         solofsapitypes.NetINodeID
		uNetINode          solofsapitypes.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
		respBody           []byte
		readDataSize       int
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.solodn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solodn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		if err == solofsapitypes.ErrObjectNotExists {
			solofsapi.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_404, "")
			goto SERVICE_REQUEST_DONE
		} else {
			solofsapi.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_502, "")
			goto SERVICE_REQUEST_DONE
		}
	}

	// calculate file data size
	if reqParam.Offset()+uint64(reqParam.Length()) > uNetINode.Ptr().Size {
		readDataSize = int(uNetINode.Ptr().Size - reqParam.Offset())
	} else {
		readDataSize = int(reqParam.Length())
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(reqParam.Offset() / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((reqParam.Offset() + uint64(readDataSize)) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.solodn.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.solodn.netBlockDriver.ReleaseNetBlock(uNetBlock)
		if err != nil {
			solofsapi.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_502, "")
			goto SERVICE_REQUEST_DONE
		}

		if uNetBlock.Ptr().IsLocalDataBackendInited.Load() == solodbapitypes.MetaDataStateUninited {
			p.solodn.PrepareNetBlockLocalDataBackend(uNetBlock)
		}
	}

SERVICE_REQUEST_DONE:
	if err != nil {
		return nil
	}

	// request file data
	solofsapi.SetNetINodePReadResponse(&protocolBuilder, int32(readDataSize))
	respBody = protocolBuilder.Bytes[protocolBuilder.Head():]
	// TODO set write length
	err = serviceReq.ResponseHeaderParam(serviceReq.ReqID, respBody, int(readDataSize))
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

	// TODO get readedDataLength
	_, err = p.solodn.netINodeDriver.PReadWithNetQuery(uNetINode, serviceReq,
		int(readDataSize), reqParam.Offset())
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

SERVICE_RESPONSE_DONE:
	return nil
}
