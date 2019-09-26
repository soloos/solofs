package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/util"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	uploadMemBlockJobChan chan solofsapitypes.UploadMemBlockJobUintptr
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver) error {
	p.driver = driver

	p.uploadMemBlockJobChan = make(chan solofsapitypes.UploadMemBlockJobUintptr, 2048)

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}
