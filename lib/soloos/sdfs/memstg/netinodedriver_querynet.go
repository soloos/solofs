package memstg

import "soloos/sdfs/types"

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode types.NetINodeUintptr) error {
	var err error

	err = p.helper.NameNodeClient.GetNetINodeMetaData(uNetINode, 0, 0, 0)
	if err != nil {
		goto PREPARE_DONE
	}

PREPARE_DONE:
	if err == nil {
		uNetINode.Ptr().IsDBMetaDataInited = true
	}
	return nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var err error

	err = p.helper.NameNodeClient.MustGetNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		goto PREPARE_DONE
	}

PREPARE_DONE:
	if err == nil {
		uNetINode.Ptr().IsDBMetaDataInited = true
	}
	return nil
}
