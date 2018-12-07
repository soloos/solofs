package namenode

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) NetBlockPrepareMetadata(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param   = make([]byte, reqBodySize)
		req     protocol.INodeNetBlockInfoRequest
		uINode  types.INodeUintptr
		inodeID types.INodeID
		err     error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(inodeID[:], req.InodeID())
	uINode, err = p.nameNode.MetaStg.GetINode(inodeID)
	if err != nil {
		goto SERVICE_DONE
	}

	// response
	util.Ignore(uINode)
	// api.SetINodeNetBlockInfoResponse(p.dataNodePeers[:], req.Cap(), req.Cap(), &protocolBuilder)

SERVICE_DONE:
	return nil
}
