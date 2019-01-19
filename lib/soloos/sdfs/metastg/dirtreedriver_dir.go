package metastg

import (
	"path/filepath"
	"soloos/sdfs/types"
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeDriver) ListFsINodeByParentPath(parentPath string,
	beforeLiteralFunc func(resultCount int) bool,
	literalFunc func(types.FsINode) bool,
) error {
	var (
		fsINode types.FsINode
		err     error
	)

	fsINode, err = p.GetFsINodeByPath(parentPath)
	if err != nil {
		return err
	}

	err = p.ListFsINodeByParentIDFromDB(fsINode.Ino, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) Rename(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINode                       types.FsINode
		oldFsINode                    types.FsINode
		parentFsINode                 types.FsINode
		tmpFsINode                    types.FsINode
		tmpParentDirPath, tmpFileName string
		err                           error
	)

	oldFsINode, err = p.GetFsINodeByPath(oldFsINodeName)
	if err != nil {
		return err
	}
	fsINode = oldFsINode

	tmpParentDirPath, tmpFileName = filepath.Split(newFsINodePath)
	parentFsINode, err = p.GetFsINodeByPath(tmpParentDirPath)
	if err != nil {
		return err
	}

	if parentFsINode.Type != types.FSINODE_TYPE_IFDIR {
		return types.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINode.ParentID = parentFsINode.Ino
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	tmpFsINode, err = p.GetFsINodeByPath(newFsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			fsINode.ParentID = parentFsINode.Ino
			fsINode.Name = tmpFileName
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return types.ErrObjectNotExists
		}
	}

	if tmpFsINode.Type == types.FSINODE_TYPE_IFDIR {
		parentFsINode = tmpFsINode
		fsINode.ParentID = parentFsINode.Ino
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return types.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	p.deleteFsINodeCache(oldFsINode.ParentID, oldFsINode.Name, oldFsINode.Ino)

	err = p.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) Mkdir(ino types.FsINodeID, input *fuse.MkdirIn, name string, output *fuse.EntryOut) error {
	var (
		fsINode  types.FsINode
		parentID types.FsINodeID = input.InHeader.NodeId
		err      error
	)

	_, err = p.GetFsINodeByName(parentID, name)
	if err == nil {
		return nil
	}

	if err != nil && err != types.ErrObjectNotExists {
		return err
	}

	now := time.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
	fsINode = types.FsINode{
		Ino:        ino,
		NetINodeID: types.ZeroNetINodeID,
		ParentID:   parentID,
		Name:       name,
		Type:       types.FSINODE_TYPE_IFDIR,
		Atime:      nowt,
		Ctime:      nowt,
		Mtime:      nowt,
		Atimensec:  nowtnsec,
		Ctimensec:  nowtnsec,
		Mtimensec:  nowtnsec,
		Mode:       input.Mode,
		Nlink:      1,
	}
	err = p.InsertFsINodeInDB(fsINode)
	if err != nil {
		return err
	}

	return err
}
