package types

import (
	sdbapitypes "soloos/common/sdbapi/types"
)

type MetaDataState = sdbapitypes.MetaDataState

const (
	MetaDataStateUninited = sdbapitypes.MetaDataStateUninited
	MetaDataStateInited   = sdbapitypes.MetaDataStateInited
)
