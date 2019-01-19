package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeDriver) MakeFsINodeKey(parentID types.FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + fsINodeName
}

func (p *DirTreeDriver) PrepareSchema() error {
	var err error
	switch p.helper.DBConn.DBDriver {
	case "mysql":
		err = InstallMysqlSchema(p.helper.DBConn)
	case "sqlite3":
		err = InstallSqlite3Schema(p.helper.DBConn)
	}
	return err
}

func (p *DirTreeDriver) PrepareINodes() error {
	var err error

	p.allocINodeIDDalta = 10000 * 10
	for p.lastFsINodeIDInDB <= types.RootFsINodeID {
		p.lastFsINodeIDInDB, err = FetchAndUpdateMaxID(p.helper.DBConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeID = p.lastFsINodeIDInDB
	}

	p.fsINodesByPath = make(map[string]types.FsINode)
	p.fsINodesByID = make(map[types.FsINodeID]types.FsINode)

	p.Mkdir(types.RootFsINodeID, &fuse.MkdirIn{
		InHeader: fuse.InHeader{
			Length: 0,
			Opcode: 0,
			Unique: 0,
			NodeId: types.RootFsINodeParentID,
			Context: fuse.Context{
				Owner: fuse.Owner{
					Uid: 0,
					Gid: 0,
				},
				Pid: 0,
			},
			Padding: 0,
		},
		Mode:  fuse.S_IFDIR | 0777,
		Umask: 0,
	}, "", &fuse.EntryOut{})

	p.Mkdir(p.AllocFsINodeID(), &fuse.MkdirIn{
		InHeader: fuse.InHeader{
			Length: 0,
			Opcode: 0,
			Unique: 0,
			NodeId: types.RootFsINodeID,
			Context: fuse.Context{
				Owner: fuse.Owner{
					Uid: 0,
					Gid: 0,
				},
				Pid: 0,
			},
			Padding: 0,
		},
		Mode:  fuse.S_IFDIR | 0777,
		Umask: 0,
	}, "tmp", &fuse.EntryOut{})

	p.rootFsINode, err = p.GetFsINodeByName(types.RootFsINodeParentID, "")
	if err != nil {
		return err
	}

	for i := 0; i < len(p.sysFsINode); i++ {
		p.sysFsINode[i].Ino = types.FsINodeID(i)
	}

	return nil
}

func (p *DirTreeDriver) ensureFsINodeHasNetINode(fsINode *types.FsINode) error {
	if fsINode.Type != types.FSINODE_TYPE_IFREG {
		return nil
	}

	if fsINode.UNetINode != 0 {
		return nil
	}

	var err error
	fsINode.UNetINode, err = p.helper.GetNetINodeWithReadAcquire(fsINode.NetINodeID)
	return err
}

func (p *DirTreeDriver) prepareAndSetFsINodeCache(fsINode *types.FsINode) error {
	var err error
	err = p.ensureFsINodeHasNetINode(fsINode)
	if err != nil {
		return err
	}

	p.fsINodesByPathRWMutex.Lock()
	p.fsINodesByPath[p.MakeFsINodeKey(fsINode.ParentID, fsINode.Name)] = *fsINode
	p.fsINodesByPathRWMutex.Unlock()

	p.fsINodesByIDRWMutex.Lock()
	p.fsINodesByID[fsINode.Ino] = *fsINode
	p.fsINodesByIDRWMutex.Unlock()

	return nil
}

func (p *DirTreeDriver) deleteFsINodeCache(parentID types.FsINodeID, fsINodeName string, fsINodeID types.FsINodeID) {
	p.fsINodesByPathRWMutex.Lock()
	delete(p.fsINodesByPath, p.MakeFsINodeKey(parentID, fsINodeName))
	p.fsINodesByPathRWMutex.Unlock()

	p.fsINodesByIDRWMutex.Lock()
	delete(p.fsINodesByID, fsINodeID)
	p.fsINodesByIDRWMutex.Unlock()
}

func (p *DirTreeDriver) AllocFsINodeID() types.FsINodeID {
	var ret = atomic.AddUint64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.helper.DBConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}

func (p *DirTreeDriver) DeleteINodeByPath(fsINodePath string) error {
	var (
		fsINode types.FsINode
		err     error
	)

	fsINode, err = p.GetFsINodeByPath(fsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.DeleteFsINodeByIDInDB(fsINode.Ino)
	p.deleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.Ino)

	return err
}

func (p *DirTreeDriver) GetFsINodeByPath(fsInodePath string) (types.FsINode, error) {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = p.rootFsINode.Ino
		fsINode  types.FsINode
		err      error
	)

	paths = strings.Split(fsInodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	if len(paths) <= 1 {
		return p.rootFsINode, nil
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByName(parentID, paths[i])
		if err != nil {
			return fsINode, err
		}
		parentID = fsINode.Ino
	}

	return fsINode, err
}

func (p *DirTreeDriver) GetFsINodeByID(fsINodeID types.FsINodeID) (types.FsINode, error) {
	var (
		fsINode types.FsINode
		exists  bool
		err     error
	)

	p.fsINodesByIDRWMutex.RLock()
	fsINode, exists = p.fsINodesByID[fsINodeID]
	p.fsINodesByIDRWMutex.RUnlock()
	if exists {
		return fsINode, nil
	}

	if fsINodeID < types.RootFsINodeID {
		return p.sysFsINode[fsINodeID], nil
	}

	fsINode, err = p.GetFsINodeByIDFromDB(fsINodeID)
	if err != nil {
		return fsINode, err
	}

	err = p.prepareAndSetFsINodeCache(&fsINode)

	return fsINode, err
}

func (p *DirTreeDriver) GetFsINodeByName(parentID types.FsINodeID, fsINodeName string) (types.FsINode, error) {
	var (
		fsINode types.FsINode
		exists  bool
		err     error
	)

	p.fsINodesByPathRWMutex.RLock()
	fsINode, exists = p.fsINodesByPath[p.MakeFsINodeKey(parentID, fsINodeName)]
	p.fsINodesByPathRWMutex.RUnlock()
	if exists {
		return fsINode, nil
	}

	fsINode, err = p.GetFsINodeByNameFromDB(parentID, fsINodeName)
	if err != nil {
		return fsINode, err
	}

	err = p.prepareAndSetFsINodeCache(&fsINode)

	return fsINode, err
}
