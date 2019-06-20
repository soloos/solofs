package memstg

import (
	"os"
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/sdfs/types"
	"strings"
)

func (p *PosixFS) SimpleOpenFile(fsINodePath string,
	netBlockCap int, memBlockCap int) (sdfsapitypes.FsINodeMeta, error) {
	var (
		paths       []string
		i           int
		parentID    sdfsapitypes.FsINodeID = sdfsapitypes.RootFsINodeID
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths)-1; i++ {
		if paths[i] == "" {
			continue
		}
		err = p.FetchFsINodeByName(&fsINodeMeta, parentID, paths[i])
		if err != nil {
			goto OPEN_FILE_DONE
		}
		parentID = fsINodeMeta.Ino
	}

	err = p.FetchFsINodeByName(&fsINodeMeta, parentID, paths[i])
	if err == nil {
		goto OPEN_FILE_DONE
	}

	if err == sdfsapitypes.ErrObjectNotExists {
		err = p.createFsINode(&fsINodeMeta,
			nil, nil, parentID,
			paths[i], types.FSINODE_TYPE_FILE, fsapitypes.S_IFREG|0777,
			0, 0, types.FS_RDEV)
		if err != nil {
			goto OPEN_FILE_DONE
		}
	}

OPEN_FILE_DONE:
	return fsINodeMeta, err
}

func (p *PosixFS) Create(input *fsapitypes.CreateIn, name string, out *fsapitypes.CreateOut) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, input.NodeId,
		name, types.FSINODE_TYPE_FILE,
		uint32(0777)&input.Mode|uint32(fsapitypes.S_IFREG),
		input.Uid, input.Gid, types.FS_RDEV)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, &out.OpenOut)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(&out.EntryOut, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFS) Open(input *fsapitypes.OpenIn, out *fsapitypes.OpenOut) fsapitypes.Status {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(input.NodeId)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&uFsINode.Ptr().Meta, input.Flags, out)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	openFlags := int(input.Flags)
	if (openFlags&os.O_TRUNC != 0) ||
		(openFlags&os.O_WRONLY != 0) ||
		(openFlags&os.O_APPEND != 0) {
		err = p.FsINodeDriver.RefreshFsINodeACMtime(uFsINode)
		if err != nil {
			return types.ErrorToFsStatus(err)
		}
	}

	return fsapitypes.OK
}

func (p *PosixFS) Fallocate(input *fsapitypes.FallocateIn) fsapitypes.Status {
	// TODO maybe should support
	// not support
	return fsapitypes.ENODATA
}
