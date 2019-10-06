package memstg

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/solomqapi"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type PrepareNetBlockMetaData func(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error

type NetBlockDriverHelper struct {
	PrepareNetBlockMetaData
}

type NetBlockDriver struct {
	*soloosbase.SoloosEnv
	helper NetBlockDriverHelper

	solonnClient *solofsapi.SolonnClient
	solodnClient *solofsapi.SolodnClient
	solomqClient solomqapi.Client

	netBlockTable          offheap.LKVTableWithBytes68
	netBlockDriverUploader netBlockDriverUploader

	preadMemBlockWithDisk    solofsapitypes.PReadMemBlockWithDisk
	uploadMemBlockWithDisk   solofsapitypes.UploadMemBlockWithDisk
	uploadMemBlockWithSolomq solofsapitypes.UploadMemBlockWithSolomq
}

func (p *NetBlockDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) error {
	var err error

	p.SoloosEnv = soloosEnv

	p.SetHelper(prepareNetBlockMetaData)

	p.solonnClient = solonnClient
	p.solodnClient = solodnClient

	err = p.OffheapDriver.InitLKVTableWithBytes68(&p.netBlockTable, "NetBlock",
		int(solofsapitypes.NetBlockStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	err = p.netBlockDriverUploader.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) netBlockTablePrepareNewObjectFunc(uNetBlock solofsapitypes.NetBlockUintptr,
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
	prepareNetBlockMetaData PrepareNetBlockMetaData,
) {
	p.helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
}

func (p *NetBlockDriver) SetSolomqClient(solomqClient solomqapi.Client) {
	p.solomqClient = solomqClient
}

func (p *NetBlockDriver) SetSolonnClient(solonnClient *solofsapi.SolonnClient) {
	p.solonnClient = solonnClient
}

func (p *NetBlockDriver) SetPReadMemBlockWithDisk(preadMemBlockWithDisk solofsapitypes.PReadMemBlockWithDisk) {
	p.preadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *NetBlockDriver) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk) {
	p.uploadMemBlockWithDisk = uploadMemBlockWithDisk
}

func (p *NetBlockDriver) SetUploadMemBlockWithSolomq(uploadMemBlockWithSolomq solofsapitypes.UploadMemBlockWithSolomq) {
	p.uploadMemBlockWithSolomq = uploadMemBlockWithSolomq
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode solofsapitypes.NetINodeUintptr,
	netBlockIndex int32) (solofsapitypes.NetBlockUintptr, error) {
	var (
		uNetBlock       solofsapitypes.NetBlockUintptr
		pNetBlock       *solofsapitypes.NetBlock
		uObject         offheap.LKVTableObjectUPtrWithBytes68
		netINodeBlockID solofsapitypes.NetINodeBlockID
		afterSetNewObj  offheap.KVTableAfterSetNewObj
		err             error
	)

	solofsapitypes.EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	uObject, afterSetNewObj = p.netBlockTable.MustGetObject(netINodeBlockID)
	p.netBlockTablePrepareNewObjectFunc(solofsapitypes.NetBlockUintptr(uObject), afterSetNewObj)
	uNetBlock = solofsapitypes.NetBlockUintptr(uObject)
	pNetBlock = uNetBlock.Ptr()
	if uNetBlock.Ptr().IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
		pNetBlock.IsDBMetaDataInited.LockContext()
		if pNetBlock.IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
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

func (p *NetBlockDriver) ReleaseNetBlock(uNetBlock solofsapitypes.NetBlockUintptr) {
	p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
}
