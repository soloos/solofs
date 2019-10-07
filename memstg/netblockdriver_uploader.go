package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/util"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobChan chan solofstypes.UploadMemBlockJobUintptr
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan solofstypes.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
