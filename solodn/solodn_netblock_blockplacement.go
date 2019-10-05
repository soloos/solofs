package solodn

import (
	"soloos/common/snettypes"
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
)

func (p *Solodn) doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock solofsapitypes.NetBlockUintptr,
	backends snettypes.PeerGroup,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsSyncDataBackendsInited.LockContext()
	if pNetBlock.IsSyncDataBackendsInited.Load() == solodbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	// fanout
	pNetBlock.SyncDataBackends.Reset()
	for i, _ := range backends.Slice() {
		pNetBlock.SyncDataBackends.Append(backends.Arr[i], 0)
	}
	pNetBlock.IsSyncDataBackendsInited.Store(solodbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsSyncDataBackendsInited.UnlockContext()
	return err
}

func (p *Solodn) PrepareNetBlockSyncDataBackends(uNetBlock solofsapitypes.NetBlockUintptr,
	backends snettypes.PeerGroup,
) error {
	return p.doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock, backends)
}

func (p *Solodn) PrepareNetBlockLocalDataBackend(uNetBlock solofsapitypes.NetBlockUintptr) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsLocalDataBackendInited.LockContext()
	if pNetBlock.IsLocalDataBackendInited.Load() == solodbapitypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	// TODO: maybe pNetBlock.IsLocalDataBackendExists is false
	pNetBlock.IsLocalDataBackendExists = true
	pNetBlock.IsLocalDataBackendInited.Store(solodbapitypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited.UnlockContext()
	return err
}
