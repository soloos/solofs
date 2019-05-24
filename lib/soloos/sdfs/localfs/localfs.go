package localfs

type LocalFS struct {
	fdDriver FdDriver
}

func (p *LocalFS) Init(dataPathPrefix string) error {
	var (
		err error
	)

	err = p.fdDriver.Init(dataPathPrefix)
	if err != nil {
		return err
	}

	return nil
}
