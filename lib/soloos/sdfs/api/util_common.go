package api

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
	"golang.org/x/xerrors"
)

func SetCommonResponseCode(protocolBuilder *flatbuffers.Builder, code int) {
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocolBuilder.Finish(protocol.CommonResponseEnd(protocolBuilder))
}

func CommonResponseToError(obj *protocol.CommonResponse) error {
	switch obj.Code() {
	case snettypes.CODE_OK:
		return nil
	case snettypes.CODE_404:
		return types.ErrObjectNotExists
	case snettypes.CODE_502:
		return types.ErrRemoteService
	}

	return xerrors.New(string(obj.Error()))
	// return types.ErrRemoteService
}
