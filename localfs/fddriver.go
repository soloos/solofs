package localfs

import (
	"os"
	"path/filepath"
	"soloos/common/solofstypes"
	"soloos/common/util"
)

type FdDriver struct {
	dataPathPrefix string
	fdMutex        util.Mutex
	fds            map[solofstypes.NetINodeUintptr]*Fd
}

func (p *FdDriver) Init(dataPathPrefix string) error {
	var err error
	p.dataPathPrefix = dataPathPrefix
	err = os.MkdirAll(p.dataPathPrefix, 0755)
	if err != nil {
		return err
	}

	p.fds = make(map[solofstypes.NetINodeUintptr]*Fd)

	return nil
}

func (p *FdDriver) Open(uNetINode solofstypes.NetINodeUintptr, uNetBlock solofstypes.NetBlockUintptr) (*Fd, error) {
	var (
		fd  *Fd
		err error
	)
	p.fdMutex.Lock()
	fd = p.fds[uNetINode]
	if fd == nil {
		fd = new(Fd)
		err = fd.Init(uNetINode, filepath.Join(p.dataPathPrefix, uNetINode.Ptr().IDStr()))
		if err != nil {
			goto OPEN_DONE
		}
		p.fds[uNetINode] = fd
	}

	fd.BorrowResource()

OPEN_DONE:
	p.fdMutex.Unlock()
	return fd, nil
}

func (p *FdDriver) Close(fd *Fd) error {
	if fd == nil {
		return nil
	}

	var (
		err error
	)

	p.fdMutex.Lock()
	fd = p.fds[fd.uNetINode]
	if fd == nil {
		goto CLOSE_DONE
	}

	fd.ReturnResource()

CLOSE_DONE:
	p.fdMutex.Unlock()
	return err
}
