package memstg

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/swalapi"
	"soloos/sdbone/offheap"
)

type PrepareNetBlockMetaData func(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error

type NetBlockDriverHelper struct {
	NameNodeClient *sdfsapi.NameNodeClient
	SWALClient     swalapi.Client
	PrepareNetBlockMetaData
}

type NetBlockDriver struct {
	*soloosbase.SoloOSEnv
	helper NetBlockDriverHelper

	netBlockTable          offheap.LKVTableWithBytes68
	dataNodeClient         *sdfsapi.DataNodeClient
	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.SetHelper(nameNodeClient, prepareNetBlockMetaData)

	err = p.OffheapDriver.InitLKVTableWithBytes68(&p.netBlockTable, "NetBlock",
		int(sdfsapitypes.NetBlockStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
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

func (p *NetBlockDriver) netBlockTablePrepareNewObjectFunc(uNetBlock sdfsapitypes.NetBlockUintptr,
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
	nameNodeClient *sdfsapi.NameNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
}

func (p *NetBlockDriver) SetSWALClient(swalClient swalapi.Client) {
	p.helper.SWALClient = swalClient
}

func (p *NetBlockDriver) SetPReadMemBlockWithDisk(preadWithDisk sdfsapitypes.PReadMemBlockWithDisk) {
	p.dataNodeClient.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *NetBlockDriver) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk sdfsapitypes.UploadMemBlockWithDisk) {
	p.dataNodeClient.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}

func (p *NetBlockDriver) SetUploadMemBlockWithSWAL(uploadMemBlockWithSWAL sdfsapitypes.UploadMemBlockWithSWAL) {
	p.dataNodeClient.SetUploadMemBlockWithSWAL(uploadMemBlockWithSWAL)
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode sdfsapitypes.NetINodeUintptr,
	netBlockIndex int32) (sdfsapitypes.NetBlockUintptr, error) {
	var (
		uNetBlock       sdfsapitypes.NetBlockUintptr
		pNetBlock       *sdfsapitypes.NetBlock
		uObject         offheap.LKVTableObjectUPtrWithBytes68
		netINodeBlockID sdfsapitypes.NetINodeBlockID
		afterSetNewObj  offheap.KVTableAfterSetNewObj
		err             error
	)

	sdfsapitypes.EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	uObject, afterSetNewObj = p.netBlockTable.MustGetObject(netINodeBlockID)
	p.netBlockTablePrepareNewObjectFunc(sdfsapitypes.NetBlockUintptr(uObject), afterSetNewObj)
	uNetBlock = sdfsapitypes.NetBlockUintptr(uObject)
	pNetBlock = uNetBlock.Ptr()
	if uNetBlock.Ptr().IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
		pNetBlock.IsDBMetaDataInited.LockContext()
		if pNetBlock.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
			err = p.helper.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
		}
		pNetBlock.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
		return 0, err
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) ReleaseNetBlock(uNetBlock sdfsapitypes.NetBlockUintptr) {
	p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
}
