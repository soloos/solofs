package sfuse

import "github.com/hanwen/go-fuse/fuse"

// Extended attributes.
func (p *SFuseFs) GetXAttrSize(header *fuse.InHeader, attr string) (sz int, code fuse.Status) {
	return p.Client.MemDirTreeStg.GetXAttrSize(header.NodeId, attr)
}

func (p *SFuseFs) GetXAttrData(header *fuse.InHeader, attr string) (data []byte, code fuse.Status) {
	return p.Client.MemDirTreeStg.GetXAttrData(header.NodeId, attr)
}

func (p *SFuseFs) ListXAttr(header *fuse.InHeader) (attributes []byte, code fuse.Status) {
	return p.Client.MemDirTreeStg.ListXAttr(header.NodeId)
}

func (p *SFuseFs) SetXAttr(input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return p.Client.MemDirTreeStg.SetXAttr(input.NodeId, attr, data)
}

func (p *SFuseFs) RemoveXAttr(header *fuse.InHeader, attr string) (code fuse.Status) {
	return p.Client.MemDirTreeStg.RemoveXAttr(header.NodeId, attr)
}
