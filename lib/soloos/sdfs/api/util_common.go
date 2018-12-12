package api

import (
	"errors"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

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
