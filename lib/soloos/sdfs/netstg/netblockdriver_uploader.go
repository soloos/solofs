package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util"
	"soloos/util/offheap"
	"sync"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver

	uploadJobMutex      sync.Mutex
	uploadJobChan       chan UploadJobUintptr
	uploadJobs          map[types.MemBlockUintptr]UploadJobUintptr
	uploadChunkMaskPool offheap.RawObjectPool
	uploadJobPool       offheap.RawObjectPool
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver) error {
	var err error
	p.driver = driver

	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver

	p.uploadJobChan = make(chan UploadJobUintptr, 2048)
	p.uploadJobs = make(map[types.MemBlockUintptr]UploadJobUintptr, 2048)
	err = p.driver.offheapDriver.InitRawObjectPool(&p.uploadChunkMaskPool,
		int(offheap.ChunkMaskStructSize), -1, nil, nil)
	if err != nil {
		return err
	}
	err = p.driver.offheapDriver.InitRawObjectPool(&p.uploadJobPool,
		int(UploadJobStructSize), -1, p.RawChunkPoolInvokePrepareNewUploadJob, nil)
	if err != nil {
		return err
	}

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
