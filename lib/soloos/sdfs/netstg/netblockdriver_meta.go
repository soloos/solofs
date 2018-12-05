package netstg

import "soloos/sdfs/types"

func (p *NetBlockDriver) PrepareNetBlockMetadata(uINode types.INodeUintptr, uNetBlock types.NetBlockUintptr) {
	pNetBlock := uNetBlock.Ptr()
	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	// p.driver.snetClientDriver.Call()

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
}
