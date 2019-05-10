package netstg

import (
	sdbapitypes "soloos/common/sdbapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
)

type PrepareNetBlockMetaData func(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int32) error

type NetBlockDriverHelper struct {
	*api.NameNodeClient
	PrepareNetBlockMetaData
}

type NetBlockDriver struct {
	*soloosbase.SoloOSEnv
	helper NetBlockDriverHelper

	netBlockTable          offheap.LKVTableWithBytes68
	dataNodeClient         *api.DataNodeClient
	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.SetHelper(nameNodeClient, prepareNetBlockMetaData)

	err = p.OffheapDriver.InitLKVTableWithBytes68(&p.netBlockTable, "NetBlock",
		int(types.NetBlockStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	p.dataNodeClient = dataNodeClient

	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) netBlockTablePrepareNewObjectFunc(uNetBlock types.NetBlockUintptr,
	afterSetNewObj offheap.KVTableAfterSetNewObj) bool {
	var isNewObjectSetted bool
	if afterSetNewObj != nil {
		afterSetNewObj()
		isNewObjectSetted = true
	} else {
		isNewObjectSetted = false
	}
	return isNewObjectSetted

}

func (p *NetBlockDriver) SetHelper(
	nameNodeClient *api.NameNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
}

func (p *NetBlockDriver) SetPReadMemBlockWithDisk(preadWithDisk api.PReadMemBlockWithDisk) {
	p.dataNodeClient.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *NetBlockDriver) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk api.UploadMemBlockWithDisk) {
	p.dataNodeClient.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}

func (p *NetBlockDriver) SyncMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr) error {
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()
	return nil
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode types.NetINodeUintptr,
	netBlockIndex int32) (types.NetBlockUintptr, error) {
	var (
		uNetBlock       types.NetBlockUintptr
		pNetBlock       *types.NetBlock
		uObject         offheap.LKVTableObjectUPtrWithBytes68
		netINodeBlockID types.NetINodeBlockID
		afterSetNewObj  offheap.KVTableAfterSetNewObj
		err             error
	)

	sdfsapitypes.EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	uObject, afterSetNewObj = p.netBlockTable.MustGetObject(netINodeBlockID)
	p.netBlockTablePrepareNewObjectFunc(types.NetBlockUintptr(uObject), afterSetNewObj)
	uNetBlock = types.NetBlockUintptr(uObject)
	pNetBlock = uNetBlock.Ptr()
	if uNetBlock.Ptr().IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
		pNetBlock.IsDBMetaDataInited.LockContext()
		if pNetBlock.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
			err = p.helper.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
		}
		pNetBlock.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netBlockTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
		return 0, err
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) ReleaseNetBlock(uNetBlock types.NetBlockUintptr) {
	p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
}
