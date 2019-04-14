package netstg

import (
	"soloos/common/snet"
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
	helper NetBlockDriverHelper

	offheapDriver *offheap.OffheapDriver
	netBlockTable offheap.LKVTableWithBytes68

	snetDriver       *snet.NetDriver
	snetClientDriver *snet.ClientDriver
	dataNodeClient   *api.DataNodeClient

	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.offheapDriver.InitLKVTableWithBytes68(&p.netBlockTable, "NetBlock",
		int(types.NetBlockStructSize), -1, types.DefaultKVTableSharedCount,
		p.netBlockTablePrepareNewObjectFunc, nil)
	if err != nil {
		return err
	}

	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.dataNodeClient = dataNodeClient

	p.SetHelper(nameNodeClient, prepareNetBlockMetaData)

	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) netBlockTablePrepareNewObjectFunc(uObject uintptr) {
	types.NetBlockUintptr(uObject).Ptr().ID = types.NetBlockUintptr(uObject).Ptr().LKVTableObjectWithBytes68.ID
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

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode types.NetINodeUintptr,
	netBlockIndex int32) (types.NetBlockUintptr, error) {
	var (
		uNetBlock       types.NetBlockUintptr
		pNetBlock       *types.NetBlock
		uObject         uintptr
		netINodeBlockID types.NetINodeBlockID
		isLoaded        bool
		err             error
	)

	types.EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	uObject, isLoaded = p.netBlockTable.MustGetObjectWithAcquire(netINodeBlockID)
	uNetBlock = types.NetBlockUintptr(uObject)
	pNetBlock = uNetBlock.Ptr()
	if isLoaded == false || uNetBlock.Ptr().IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
		pNetBlock.DBMetaDataInitMutex.Lock()
		if pNetBlock.IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
			err = p.helper.PrepareNetBlockMetaData(uNetBlock, uNetINode, netBlockIndex)
		}
		pNetBlock.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
		return 0, err
	}

	return uNetBlock, nil
}

func (p *NetBlockDriver) FlushMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr) error {
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()
	return nil
}

func (p *NetBlockDriver) GetDataNodeClient() *api.DataNodeClient {
	return p.dataNodeClient
}
