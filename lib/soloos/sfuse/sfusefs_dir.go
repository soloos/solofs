package sfuse

import (
	"soloos/log"
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *SFuseFs) Mkdir(input *fuse.MkdirIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	log.Error("fuck you shit")
	var err error
	err = p.Client.DirTreeDriver.Mkdir(p.Client.DirTreeDriver.AllocFsINodeID(), input, name, out)
	if err == types.ErrObjectExists {
		return fuse.EPERM
	}
	return fuse.EPERM
}

func (p *SFuseFs) Rmdir(header *fuse.InHeader, name string) (code fuse.Status) {
	log.Error("fuck you shit")
	return fuse.EPERM
}

// Directory handling
func (p *SFuseFs) OpenDir(input *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	log.Error("fuck you shit")
	return fuse.EPERM
}

func (p *SFuseFs) ReadDir(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	log.Error("fuck you shit")
	return fuse.EPERM
}

func (p *SFuseFs) ReadDirPlus(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	log.Error("fuck you shit")
	return fuse.EPERM
}

func (p *SFuseFs) ReleaseDir(input *fuse.ReleaseIn) {
	log.Error("fuck you shit")
}

func (p *SFuseFs) FsyncDir(input *fuse.FsyncIn) (code fuse.Status) {
	log.Error("fuck you shit")
	return fuse.EPERM
}
