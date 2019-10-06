package memstg

import (
	"os"
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
	"strings"
)

func (p *PosixFs) SimpleOpenFile(fsINodePath string,
	netBlockCap int, memBlockCap int) (solofsapitypes.FsINodeMeta, error) {
	var (
		paths       []string
		i           int
		parentID    solofsapitypes.FsINodeID = solofsapitypes.RootFsINodeID
		fsINodeMeta solofsapitypes.FsINodeMeta
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

	if err.Error() == solofsapitypes.ErrObjectNotExists.Error() {
		err = p.createFsINode(&fsINodeMeta,
			nil, nil, parentID,
			paths[i], solofstypes.FSINODE_TYPE_FILE, fsapitypes.S_IFREG|0777,
			0, 0, solofstypes.FS_RDEV)
		if err != nil {
			goto OPEN_FILE_DONE
		}
	}

OPEN_FILE_DONE:
	return fsINodeMeta, err
}

func (p *PosixFs) Create(input *fsapitypes.CreateIn, name string, out *fsapitypes.CreateOut) fsapitypes.Status {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)

	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, input.NodeId,
		name, solofstypes.FSINODE_TYPE_FILE,
		uint32(0777)&input.Mode|uint32(fsapitypes.S_IFREG),
		input.Uid, input.Gid, solofstypes.FS_RDEV)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, &out.OpenOut)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(&out.EntryOut, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFs) Open(input *fsapitypes.OpenIn, out *fsapitypes.OpenOut) fsapitypes.Status {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(input.NodeId)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&uFsINode.Ptr().Meta, input.Flags, out)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	openFlags := int(input.Flags)
	if (openFlags&os.O_TRUNC != 0) ||
		(openFlags&os.O_WRONLY != 0) ||
		(openFlags&os.O_APPEND != 0) {
		err = p.FsINodeDriver.RefreshFsINodeACMtime(uFsINode)
		if err != nil {
			return solofstypes.ErrorToFsStatus(err)
		}
	}

	return fsapitypes.OK
}

func (p *PosixFs) Fallocate(input *fsapitypes.FallocateIn) fsapitypes.Status {
	// TODO maybe should support
	// not support
	return fsapitypes.ENODATA
}
