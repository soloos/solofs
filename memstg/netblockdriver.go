package memstg

import (
	"soloos/common/solodbtypes"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"soloos/common/solomqapi"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type PrepareNetBlockMetaData func(uNetBlock solofstypes.NetBlockUintptr,
	uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error

type NetBlockDriverHelper struct {
	PrepareNetBlockMetaData
	PreadMemBlockWithDisk    solofstypes.PReadMemBlockWithDisk
	UploadMemBlockWithDisk   solofstypes.UploadMemBlockWithDisk
	UploadMemBlockWithSolomq solofstypes.UploadMemBlockWithSolomq
}

type NetBlockDriver struct {
	*soloosbase.SoloosEnv
	helper NetBlockDriverHelper

	solonnClient *solofsapi.SolonnClient
	solodnClient *solofsapi.SolodnClient
	solomqClient solomqapi.Client

	netBlockTable          offheap.LKVTableWithBytes68
	NetBlockUploader NetBlockUploader
}

func (p *NetBlockDriver) Init(
	soloosEnv *soloosbase.SoloosEnv,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	prepareNetBlockMetaData PrepareNetBlockMetaData,
	preadMemBlockWithDisk solofstypes.PReadMemBlockWithDisk,
	uploadMemBlockWithDisk solofstypes.UploadMemBlockWithDisk,
	uploadMemBlockWithSolomq solofstypes.UploadMemBlockWithSolomq,
) error {
	var err error

	p.SoloosEnv = soloosEnv

	p.SetHelper(
		prepareNetBlockMetaData,
		preadMemBlockWithDisk,
		uploadMemBlockWithDisk,
		uploadMemBlockWithSolomq,
	)

	p.solonnClient = solonnClient
	p.solodnClient = solodnClient

	err = p.OffheapDriver.InitLKVTableWithBytes68(&p.netBlockTable, "NetBlock",
		int(solofstypes.NetBlockStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	err = p.NetBlockUploader.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) netBlockTablePrepareNewObjectFunc(uNetBlock solofstypes.NetBlockUintptr,
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
	preadMemBlockWithDisk solofstypes.PReadMemBlockWithDisk,
	uploadMemBlockWithDisk solofstypes.UploadMemBlockWithDisk,
	uploadMemBlockWithSolomq solofstypes.UploadMemBlockWithSolomq,
) {
	p.helper.PrepareNetBlockMetaData = prepareNetBlockMetaData
	p.helper.PreadMemBlockWithDisk = preadMemBlockWithDisk
	p.helper.UploadMemBlockWithDisk = uploadMemBlockWithDisk
	p.helper.UploadMemBlockWithSolomq = uploadMemBlockWithSolomq
}

func (p *NetBlockDriver) SetSolomqClient(solomqClient solomqapi.Client) {
	p.solomqClient = solomqClient
}

func (p *NetBlockDriver) SetSolonnClient(solonnClient *solofsapi.SolonnClient) {
	p.solonnClient = solonnClient
}

func (p *NetBlockDriver) SetSolodnClient(solodnClient *solofsapi.SolodnClient) {
	p.solodnClient = solodnClient
}

func (p *NetBlockDriver) SetPReadMemBlockWithDisk(preadMemBlockWithDisk solofstypes.PReadMemBlockWithDisk) {
	p.helper.PreadMemBlockWithDisk = preadMemBlockWithDisk
}

func (p *NetBlockDriver) SetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofstypes.UploadMemBlockWithDisk) {
	p.helper.UploadMemBlockWithDisk = uploadMemBlockWithDisk
}

func (p *NetBlockDriver) SetUploadMemBlockWithSolomq(uploadMemBlockWithSolomq solofstypes.UploadMemBlockWithSolomq) {
	p.helper.UploadMemBlockWithSolomq = uploadMemBlockWithSolomq
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockDriver) MustGetNetBlock(uNetINode solofstypes.NetINodeUintptr,
	netBlockIndex solofstypes.NetBlockIndex) (solofstypes.NetBlockUintptr, error) {
	var (
		uNetBlock       solofstypes.NetBlockUintptr
		pNetBlock       *solofstypes.NetBlock
		uObject         offheap.LKVTableObjectUPtrWithBytes68
		netINodeBlockID solofstypes.NetINodeBlockID
		afterSetNewObj  offheap.KVTableAfterSetNewObj
		err             error
	)

	solofstypes.EncodeNetINodeBlockID(&netINodeBlockID, uNetINode.Ptr().ID, netBlockIndex)
	uObject, afterSetNewObj = p.netBlockTable.MustGetObject(netINodeBlockID)
	p.netBlockTablePrepareNewObjectFunc(solofstypes.NetBlockUintptr(uObject), afterSetNewObj)
	uNetBlock = solofstypes.NetBlockUintptr(uObject)
	pNetBlock = uNetBlock.Ptr()
	if uNetBlock.Ptr().IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
		pNetBlock.IsDBMetaDataInited.LockContext()
		if pNetBlock.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
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

func (p *NetBlockDriver) ReleaseNetBlock(uNetBlock solofstypes.NetBlockUintptr) {
	p.netBlockTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes68(uNetBlock))
}
