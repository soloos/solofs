package netstg

import "soloos/sdfs/types"

func (p *NetBlockDriver) PrepareNetBlockMetadata(uINode types.INodeUintptr,
	netblockIndex int,
	uNetBlock types.NetBlockUintptr,
) error {
	var err error
	pNetBlock := uNetBlock.Ptr()
	pNetBlock.MetaDataMutex.Lock()
	if pNetBlock.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.nameNodeClient.PrepareNetBlockMetadata(uINode, netblockIndex, uNetBlock)
	if err != nil {
		goto PREPARE_DONE
	}

PREPARE_DONE:
	pNetBlock.MetaDataMutex.Unlock()
	return err
}
