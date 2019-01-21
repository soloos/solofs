package memstg

import (
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeStg) OpenFile(fsINodePath string, netBlockCap int, memBlockCap int) (types.FsINode, error) {
	var (
		paths     []string
		i         int
		parentID  types.FsINodeID = types.RootFsINodeID
		fsINode   types.FsINode
		uNetINode types.NetINodeUintptr
		err       error
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths)-1; i++ {
		if paths[i] == "" {
			continue
		}
		err = p.FsINodeDriver.FetchFsINodeByName(parentID, paths[i], &fsINode)
		if err != nil {
			goto OPEN_FILE_DONE
		}
		parentID = fsINode.Ino
	}

	err = p.FsINodeDriver.FetchFsINodeByName(parentID, paths[i], &fsINode)
	if err == nil {
		goto OPEN_FILE_DONE
	}

	if err == types.ErrObjectNotExists {
		p.FsINodeDriver.AllocNetINodeID(&fsINode)
		uNetINode, err = p.FsINodeDriver.helper.MustGetNetINodeWithReadAcquire(fsINode.NetINodeID,
			0, netBlockCap, memBlockCap)
		if err != nil {
			goto OPEN_FILE_DONE
		}

		now := p.FsINodeDriver.Timer.Now()
		nowt := types.DirTreeTime(now.Unix())
		nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
		fsINode = types.FsINode{
			Ino:        p.FsINodeDriver.helper.AllocFsINodeID(),
			NetINodeID: uNetINode.Ptr().ID,
			ParentID:   parentID,
			Name:       paths[i],
			Type:       types.FSINODE_TYPE_FILE,
			Atime:      nowt,
			Ctime:      nowt,
			Mtime:      nowt,
			Atimensec:  nowtnsec,
			Ctimensec:  nowtnsec,
			Mtimensec:  nowtnsec,
			Mode:       0,
			Nlink:      1,
		}
		err = p.FsINodeDriver.helper.InsertFsINodeInDB(fsINode)
	}

OPEN_FILE_DONE:
	if err == nil {
		err = p.FsINodeDriver.PrepareAndSetFsINodeCache(&fsINode)
	}
	return fsINode, err
}
