package solodn

import (
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
)

func (p *Solodn) doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock solofstypes.NetBlockUintptr,
	backends snet.PeerGroup,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsSyncDataBackendsInited.LockContext()
	if pNetBlock.IsSyncDataBackendsInited.Load() == solodbtypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	// fanout
	pNetBlock.SyncDataBackends.Reset()
	for i, _ := range backends.Slice() {
		pNetBlock.SyncDataBackends.Append(backends.Arr[i], 0)
	}
	pNetBlock.IsSyncDataBackendsInited.Store(solodbtypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsSyncDataBackendsInited.UnlockContext()
	return err
}

func (p *Solodn) PrepareNetBlockSyncDataBackends(uNetBlock solofstypes.NetBlockUintptr,
	backends snet.PeerGroup,
) error {
	return p.doPrepareNetBlockSyncDataBackendsWithFanout(uNetBlock, backends)
}

func (p *Solodn) PrepareNetBlockLocalDataBackend(uNetBlock solofstypes.NetBlockUintptr) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)

	pNetBlock.IsLocalDataBackendInited.LockContext()
	if pNetBlock.IsLocalDataBackendInited.Load() == solodbtypes.MetaDataStateInited {
		goto PREPARE_DONE
	}

	// TODO: maybe pNetBlock.IsLocalDataBackendExists is false
	pNetBlock.IsLocalDataBackendExists = true
	pNetBlock.IsLocalDataBackendInited.Store(solodbtypes.MetaDataStateInited)

PREPARE_DONE:
	pNetBlock.IsLocalDataBackendInited.UnlockContext()
	return err
}
