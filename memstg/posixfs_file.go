package memstg

import (
	"os"
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
	"strings"
)

func (p *PosixFs) SimpleOpenFile(fsINodePath string,
	netBlockCap int, memBlockCap int) (solofstypes.FsINodeMeta, error) {
	var (
		paths       []string
		i           int
		parentID    solofstypes.FsINodeID = solofstypes.RootFsINodeID
		fsINodeMeta solofstypes.FsINodeMeta
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

	if err.Error() == solofstypes.ErrObjectNotExists.Error() {
		err = p.createFsINode(&fsINodeMeta,
			nil, nil, parentID,
			paths[i], solofstypes.FSINODE_TYPE_FILE, fsapi.S_IFREG|0777,
			0, 0, solofstypes.FS_RDEV)
		if err != nil {
			goto OPEN_FILE_DONE
		}
	}

OPEN_FILE_DONE:
	return fsINodeMeta, err
}

func (p *PosixFs) Create(input *fsapi.CreateIn, name string, out *fsapi.CreateOut) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, input.NodeId,
		name, solofstypes.FSINODE_TYPE_FILE,
		uint32(0777)&input.Mode|uint32(fsapi.S_IFREG),
		input.Uid, input.Gid, solofstypes.FS_RDEV)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, &out.OpenOut)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(&out.EntryOut, &fsINodeMeta)

	return fsapi.OK
}

func (p *PosixFs) Open(input *fsapi.OpenIn, out *fsapi.OpenOut) fsapi.Status {
	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(input.NodeId)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&uFsINode.Ptr().Meta, input.Flags, out)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	openFlags := int(input.Flags)
	if (openFlags&os.O_TRUNC != 0) ||
		(openFlags&os.O_WRONLY != 0) ||
		(openFlags&os.O_APPEND != 0) {
		err = p.FsINodeDriver.RefreshFsINodeACMtime(uFsINode)
		if err != nil {
			return ErrorToFsStatus(err)
		}
	}

	return fsapi.OK
}

func (p *PosixFs) Fallocate(input *fsapi.FallocateIn) fsapi.Status {
	// TODO maybe should support
	// not support
	return fsapi.ENODATA
}
