package netstg

import (
	"soloos/common/util"
	"soloos/sdfs/types"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobChan chan types.UploadMemBlockJobUintptr
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan types.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
