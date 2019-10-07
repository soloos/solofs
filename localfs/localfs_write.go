package localfs

import "soloos/common/solofstypes"

func (p *LocalFs) UploadMemBlockWithDisk(uJob solofstypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
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

	return nil
}
