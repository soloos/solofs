package localfs

import (
	"soloos/sdfs/types"
)

type LocalFs struct {
	fdDriver FdDriver
}

func (p *LocalFs) Init(dataPathPrefix string) error {
	var (
		err error
	)

	err = p.fdDriver.Init(dataPathPrefix)
	if err != nil {
		return err
	}

	return nil
}

func (p *LocalFs) UploadMemBlockWithDisk(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int,
) error {
	var (
		fd  *Fd
		err error
	)

	fd, err = p.fdDriver.Open(uJob.Ptr().UNetINode, uJob.Ptr().UNetBlock)
	if err != nil {
		goto UPLOAD_DONE
	}

	err = fd.Upload(uJob)
	if err != nil {
		goto UPLOAD_DONE
	}

UPLOAD_DONE:
	// TODO catch close file error
	p.fdDriver.Close(fd)

	return types.ErrObjectNotExists
}
