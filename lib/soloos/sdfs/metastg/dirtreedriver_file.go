package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"strings"
	"time"
)

func (p *DirTreeDriver) allocNetINode(netBlockCap int, memBlockCap int) (types.NetINodeID, error) {
	var (
		netINodeID types.NetINodeID
		err        error
	)

	util.InitUUID64(&netINodeID)
	_, err = p.helper.MustGetNetINodeWithReadAcquire(netINodeID, 0, netBlockCap, memBlockCap)
	return netINodeID, err
}

func (p *DirTreeDriver) OpenFile(fsInodePath string, netBlockCap int, memBlockCap int) (types.FsINode, error) {
	var (
		paths      []string
		i          int
		parentID   types.FsINodeID = p.rootFsINode.Ino
		fsINode    types.FsINode
		netINodeID types.NetINodeID
		err        error
	)

	paths = strings.Split(fsInodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths)-1; i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByName(parentID, paths[i])
		if err != nil {
			goto OPEN_FILE_DONE
		}
		parentID = fsINode.Ino
	}

	fsINode, err = p.GetFsINodeByName(parentID, paths[i])
	if err == nil {
		goto OPEN_FILE_DONE
	}

	if err == types.ErrObjectNotExists {
		netINodeID, err = p.allocNetINode(netBlockCap, memBlockCap)
		if err != nil {
			goto OPEN_FILE_DONE
		}

		now := time.Now()
		nowt := types.DirTreeTime(now.Unix())
		nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
		fsINode = types.FsINode{
			Ino:        p.AllocFsINodeID(),
			NetINodeID: netINodeID,
			ParentID:   parentID,
			Name:       paths[i],
			Type:       types.FSINODE_TYPE_IFREG,
			Atime:      nowt,
			Ctime:      nowt,
			Mtime:      nowt,
			Atimensec:  nowtnsec,
			Ctimensec:  nowtnsec,
			Mtimensec:  nowtnsec,
			Mode:       0,
			Nlink:      1,
		}
		err = p.InsertFsINodeInDB(fsINode)
	}

OPEN_FILE_DONE:
	if err == nil {
		err = p.prepareAndSetFsINodeCache(&fsINode)
	}
	return fsINode, err
}
