package netstg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"sync"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobMutex sync.Mutex
	uploadMemBlockJobChan  chan types.UploadMemBlockJobUintptr
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan types.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
