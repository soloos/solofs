package memstg

import (
	"os"
	fsapitypes "soloos/fsapi/types"
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeStg) SimpleOpenFile(fsINodePath string, netBlockCap int, memBlockCap int) (types.FsINode, error) {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = types.RootFsINodeID
		fsINode  types.FsINode
		err      error
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
		err = p.CreateFsINode(&fsINode,
			nil, nil, parentID,
			paths[i], types.FSINODE_TYPE_FILE, fsapitypes.S_IFREG|0777,
			0, 0, types.FS_RDEV)
		if err != nil {
			goto OPEN_FILE_DONE
		}

		_, err = p.FsINodeDriver.helper.MustGetNetINodeWithReadAcquire(fsINode.NetINodeID,
			0, netBlockCap, memBlockCap)
		if err != nil {
			goto OPEN_FILE_DONE
		}
	}

OPEN_FILE_DONE:
	if err == nil {
		err = p.FsINodeDriver.PrepareAndSetFsINodeCache(&fsINode)
	}
	return fsINode, err
}

func (p *DirTreeStg) Create(input *fsapitypes.CreateIn, name string, out *fsapitypes.CreateOut) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	err = p.CreateFsINode(&fsINode,
		nil, nil, input.NodeId,
		name, types.FSINODE_TYPE_FILE,
		uint32(0777)&input.Mode|uint32(fsapitypes.S_IFREG),
		input.Uid, input.Gid, types.FS_RDEV)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINode, input.Flags, &out.OpenOut)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(fsINode.ParentID)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(&out.EntryOut, &fsINode)

	return fsapitypes.OK
}

func (p *DirTreeStg) Open(input *fsapitypes.OpenIn, out *fsapitypes.OpenOut) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINode, input.Flags, out)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	openFlags := int(input.Flags)
	if (openFlags&os.O_TRUNC != 0) ||
		(openFlags&os.O_WRONLY != 0) ||
		(openFlags&os.O_APPEND != 0) {
		err = p.RefreshFsINodeACMtime(&fsINode)
		if err != nil {
			return types.ErrorToFsStatus(err)
		}
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) Fallocate(input *fsapitypes.FallocateIn) fsapitypes.Status {
	// TODO maybe should support
	// not support
	return fsapitypes.ENODATA
}
