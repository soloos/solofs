package localfs

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
