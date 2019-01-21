package sfuse

import (
	"os"
	"soloos/sdfs/libsdfs"

	"github.com/hanwen/go-fuse/fuse"
)

type SFuseFs struct {
	ClientDriver *libsdfs.ClientDriver
	Client       libsdfs.Client
}

func (p *SFuseFs) InitBySFuse(options Options, clientDriver *libsdfs.ClientDriver) error {
	var err error
	p.ClientDriver = clientDriver
	err = p.ClientDriver.InitClient(&p.Client)
	if err != nil {
		return err
	}

	os.MkdirAll(options.MountPoint, 0777)

	return nil
}

var _ = fuse.RawFileSystem(&SFuseFs{})

// This is called on processing the first request. The
// filesystem implementation can use the server argument to
// talk back to the kernel (through notify methods).
func (p *SFuseFs) Init(server *fuse.Server) {
}

func (p *SFuseFs) String() string {
	return FuseName
}

// If called, provide debug output through the log package.
func (p *SFuseFs) SetDebug(debug bool) {
}

func (p *SFuseFs) StatFs(input *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
	capacity, files := p.Client.MemDirTreeStg.StatLimits()
	usedSize, filesCount := p.Client.MemDirTreeStg.StatFs()
	blksize := p.Client.MemDirTreeStg.BlkSize()

	out.Blocks = capacity / uint64(blksize)
	out.Bfree = (capacity - usedSize) / uint64(blksize)
	out.Bavail = (capacity - usedSize) / uint64(blksize)
	out.Files = files
	out.Ffree = out.Files - filesCount
	out.Bsize = blksize
	out.NameLen = 32767
	out.Frsize = blksize

	return fuse.OK
}

func (p *SFuseFs) Close() error {
	var err error
	err = p.Client.Close()
	if err != nil {
		return err
	}

	return nil
}
