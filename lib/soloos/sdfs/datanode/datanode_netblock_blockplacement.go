package datanode

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
)

func (p *DataNode) doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock sdfsapitypes.NetBlockUintptr,
	backends snettypes.PeerGroup,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsSyncDataBackendsInited.LockContext()
	if pNetBlock.IsSyncDataBackendsInited.Load() == sdbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	// fanout
	pNetBlock.SyncDataBackends.Reset()
	for i, _ := range backends.Slice() {
		pNetBlock.SyncDataBackends.Append(backends.Arr[i], 0)
	}
	pNetBlock.IsSyncDataBackendsInited.Store(sdbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsSyncDataBackendsInited.UnlockContext()
	return err
}

func (p *DataNode) PrepareNetBlockSyncDataBackends(uNetBlock sdfsapitypes.NetBlockUintptr,
	backends snettypes.PeerGroup,
) error {
	return p.doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock, backends)
}

func (p *DataNode) PrepareNetBlockLocalDataBackend(uNetBlock sdfsapitypes.NetBlockUintptr) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsLocalDataBackendInited.LockContext()
	if pNetBlock.IsLocalDataBackendInited.Load() == sdbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	pNetBlock.IsLocalDataBackendExists = true
	pNetBlock.IsLocalDataBackendInited.Store(sdbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited.UnlockContext()
	return err
}
