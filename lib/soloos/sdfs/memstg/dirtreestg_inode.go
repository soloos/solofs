package memstg

import (
	"soloos/sdfs/types"
)

func (p *DirTreeStg) UpdateFsINodeInDB(pFsINode *types.FsINode) error {
	return p.FsINodeDriver.UpdateFsINodeInDB(pFsINode)
}

// maybe param inode means link
func (p *DirTreeStg) returnRealINode(pFsINode *types.FsINode) types.FsINodeID {
	// TODO complete me
	return pFsINode.Ino
}

func (p *DirTreeStg) FetchFsINodeByID(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	// TODO use returnRealINode maybe fsinode is link
	return p.FsINodeDriver.FetchFsINodeByID(fsINodeID, fsINode)
}

func (p *DirTreeStg) FetchFsINodeByName(parentID types.FsINodeID, fsINodeName string, fsINode *types.FsINode) error {
	// TODO use returnRealINode maybe fsinode is link
	return p.FsINodeDriver.FetchFsINodeByName(parentID, fsINodeName, fsINode)
}

func (p *DirTreeStg) CreateINode(parentID types.FsINodeID, name string, fsINodeType int, mode uint32) error {
	var (
		fsINode types.FsINode
		err     error
	)
	p.FsINodeDriver.PrepareFsINodeForCreate(&fsINode, nil, parentID, name, fsINodeType, mode)
	err = p.FsINodeDriver.CreateINode(&fsINode)
	return err
}
