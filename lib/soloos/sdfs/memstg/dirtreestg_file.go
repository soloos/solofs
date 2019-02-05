package memstg

import (
	"os"
	"soloos/sdfs/types"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) OpenFile(fsINodePath string, netBlockCap int, memBlockCap int) (types.FsINode, error) {
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
			paths[i], types.FSINODE_TYPE_FILE, fuse.S_IFREG|0777,
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

func (p *DirTreeStg) Create(input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
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
		uint32(0777)&input.Mode|uint32(fuse.S_IFREG),
		input.Uid, input.Gid, types.FS_RDEV)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.SimpleOpen(&fsINode, input.Flags, &out.OpenOut)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(fsINode.ParentID)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseEntryOutByFsINode(&out.EntryOut, &fsINode)

	return fuse.OK
}

func (p *DirTreeStg) Open(input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.SimpleOpen(&fsINode, input.Flags, out)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	openFlags := int(input.Flags)
	if (openFlags&os.O_TRUNC != 0) ||
		(openFlags&os.O_WRONLY != 0) ||
		(openFlags&os.O_APPEND != 0) {
		err = p.RefreshFsINodeACMtime(&fsINode)
		if err != nil {
			return types.ErrorToFuseStatus(err)
		}
	}

	return fuse.OK
}

func (p *DirTreeStg) Fallocate(input *fuse.FallocateIn) fuse.Status {
	// TODO maybe should support
	// not support
	return fuse.ENODATA
}

func (p *DirTreeStg) Flush(input *fuse.FlushIn) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	if fsINode.UNetINode != 0 {
		err = p.MemStg.NetINodeDriver.Flush(fsINode.UNetINode)
		if err != nil {
			return types.ErrorToFuseStatus(err)
		}
	}

	return fuse.OK
}
