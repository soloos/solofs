package api

import (
	"errors"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetCommonResponseCode(protocolBuilder *flatbuffers.Builder, code int) error {
	var err error
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocolBuilder.Finish(protocol.CommonResponseEnd(protocolBuilder))
	return err
}

func CommonResponseToError(obj *protocol.CommonResponse) error {
	switch obj.Code() {
	case snettypes.CODE_OK:
		return nil
	case snettypes.CODE_404:
		return types.ErrObjectNotExists
	}

	return errors.New(string(obj.Error()))
	// return types.ErrRemoteService
}
