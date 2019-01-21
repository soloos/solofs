package memstg

import (
	"soloos/sdfs/types"
)

func (p *DirTreeStg) Link(pointedToFsINodeID types.FsINodeID, parentID types.FsINodeID,
	filename string, fsINode *types.FsINode) error {
	return p.FsINodeDriver.Link(pointedToFsINodeID, parentID, filename, fsINode)
}

func (p *DirTreeStg) Symlink(parentID types.FsINodeID, pointedTo string, linkName string, fsINode *types.FsINode) error {
	return p.FsINodeDriver.Symlink(parentID, pointedTo, linkName, fsINode)
}

func (p *DirTreeStg) Readlink(fsINodeID types.FsINodeID) ([]byte, error) {
	return p.FsINodeDriver.Readlink(fsINodeID)
}
