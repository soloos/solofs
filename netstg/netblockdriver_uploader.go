package netstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/util"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobChan chan sdfsapitypes.UploadMemBlockJobUintptr
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan sdfsapitypes.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
