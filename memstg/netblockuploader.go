package memstg

import (
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
	"soloos/common/util"
)

type NetBlockUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobChan chan solofstypes.UploadMemBlockJobUintptr
}

func (p *NetBlockUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan solofstypes.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}

func (p *NetBlockUploader) PrepareUploadMemBlockJob(pJob *solofstypes.UploadMemBlockJob,
	uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr, memBlockIndex int32,
	backends snet.PeerGroup) {
	pJob.MetaDataState.LockContext()
	if pJob.MetaDataState.Load() == solodbtypes.MetaDataStateInited {
		pJob.MetaDataState.UnlockContext()
		return
	}
	pJob.UNetINode = uNetINode
	pJob.UNetBlock = uNetBlock
	pJob.UMemBlock = uMemBlock
	pJob.MemBlockIndex = memBlockIndex

	pJob.UploadBlockJob.Reset()

	pJob.MetaDataState.Store(solodbtypes.MetaDataStateInited)
	pJob.MetaDataState.UnlockContext()
}
