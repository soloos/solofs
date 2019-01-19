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

func (p *SFuseFs) Forget(nodeid, nlookup uint64) {
}

// Modifying structure.
func (p *SFuseFs) Mknod(input *fuse.MknodIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Unlink(header *fuse.InHeader, name string) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Rename(input *fuse.RenameIn, oldName string, newName string) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Link(input *fuse.LinkIn, filename string, out *fuse.EntryOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Symlink(header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Readlink(header *fuse.InHeader) (out []byte, code fuse.Status) {
	return nil, fuse.EPERM
}

func (p *SFuseFs) Access(input *fuse.AccessIn) (code fuse.Status) {
	return fuse.EPERM
}

// Extended attributes.
func (p *SFuseFs) GetXAttrSize(header *fuse.InHeader, attr string) (sz int, code fuse.Status) {
	return 0, fuse.EPERM
}

func (p *SFuseFs) GetXAttrData(header *fuse.InHeader, attr string) (data []byte, code fuse.Status) {
	return nil, fuse.EPERM
}

func (p *SFuseFs) ListXAttr(header *fuse.InHeader) (attributes []byte, code fuse.Status) {
	return nil, fuse.EPERM
}

func (p *SFuseFs) SetXAttr(input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return fuse.EPERM
}

func (p *SFuseFs) RemoveXAttr(header *fuse.InHeader, attr string) (code fuse.Status) {
	return fuse.EPERM
}

// File handling.
func (p *SFuseFs) Create(input *fuse.CreateIn, name string, out *fuse.CreateOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Open(input *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	return fuse.EPERM
}

// File locking
func (p *SFuseFs) GetLk(input *fuse.LkIn, out *fuse.LkOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) SetLk(input *fuse.LkIn) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) SetLkw(input *fuse.LkIn) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Release(input *fuse.ReleaseIn) {
}

func (p *SFuseFs) Flush(input *fuse.FlushIn) fuse.Status {
	return fuse.EPERM
}

func (p *SFuseFs) Fsync(input *fuse.FsyncIn) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Fallocate(input *fuse.FallocateIn) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) StatFs(input *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
	capacity, files := p.Client.DirTreeDriver.StatLimits()
	usedSize, filesCount := p.Client.DirTreeDriver.StatFs()
	blksize := p.Client.DirTreeDriver.BlkSize()

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
