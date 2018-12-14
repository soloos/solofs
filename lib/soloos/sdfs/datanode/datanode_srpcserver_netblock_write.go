package datanode

import (
	"soloos/log"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetBlockPWrite(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		reqParamData = make([]byte, reqParamSize)
		reqParam     protocol.NetBlockPWriteRequest
		backend      protocol.NetBlockBackend
		// uNetBlock    types.NetBlockUintptr
		// uMemBlock    types.MemBlockUintptr
		i   int
		err error
	)

	// request param
	err = conn.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:reqParamSize], flatbuffers.GetUOffsetT(reqParamData[:reqParamSize]))

	// uNetBlock = p.dataNode.netBlockDriver.MustGetBlock()
	log.Error(string(reqParam.NetBlockID()))
	log.Error(reqParam.Length())
	log.Error(reqParam.Offset())
	for i = 0; i < reqParam.TransferBackendsLength(); i++ {
		reqParam.TransferBackends(&backend, i)
		log.Error(string(backend.Address()))
	}

	// request file data
	var tmp = make([]byte, reqBodySize-reqParamSize)
	err = conn.ReadAll(tmp)
	if err != nil {
		return err
	}

	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = conn.SimpleResponse(reqID, respBody)
	if err != nil {
		return err
	}

	return nil
}
