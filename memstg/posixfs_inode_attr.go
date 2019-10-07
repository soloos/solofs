package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func (p *PosixFs) SetFsEntryOutByFsINode(fsEntryOut *fsapi.EntryOut,
	pFsINodeMeta *solofstypes.FsINodeMeta) {

	fsEntryOut.NodeId = pFsINodeMeta.Ino
	fsEntryOut.EntryValid = p.FsINodeDriver.EntryAttrValid
	fsEntryOut.EntryValidNsec = p.FsINodeDriver.EntryAttrValidNsec
	fsEntryOut.AttrValid = p.FsINodeDriver.EntryAttrValid
	fsEntryOut.AttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsEntryOut.Attr, pFsINodeMeta)
}

func (p *PosixFs) SetFsINodeByFsAttr(pFsINodeMeta *solofstypes.FsINodeMeta,
	input *fsapi.SetAttrIn) fsapi.Status {

	if input.Valid&fsapi.FATTR_MODE != 0 {
		// pFsINodeMeta.Mode = uint32(0777)&input.Mode | uint32(solofstypes.FsINodeTypeToFsType(pFsINodeMeta.Type))
		pFsINodeMeta.Mode = input.Mode
	}

	if input.Valid&fsapi.FATTR_UID != 0 {
		pFsINodeMeta.Uid = input.Uid
	}
	if input.Valid&fsapi.FATTR_GID != 0 {
		pFsINodeMeta.Gid = input.Gid
	}

	if input.Valid&fsapi.FATTR_SIZE != 0 {
		if pFsINodeMeta.Type == solofstypes.FSINODE_TYPE_DIR {
			return fsapi.EISDIR
		}
		var err = p.TruncateINode(pFsINodeMeta, input.Size)
		if err != nil {
			return fsapi.EIO
		}
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := solofstypes.DirTreeTime(now.Unix())
	nowtnsec := solofstypes.DirTreeTimeNsec(now.UnixNano())

	if input.Valid&(fsapi.FATTR_ATIME|fsapi.FATTR_MTIME|fsapi.FATTR_ATIME_NOW|fsapi.FATTR_MTIME_NOW) != 0 {
		if input.Valid&fsapi.FATTR_ATIME != 0 {
			if input.Valid&fsapi.FATTR_ATIME_NOW != 0 {
				pFsINodeMeta.Atime = nowt
				pFsINodeMeta.Atimensec = nowtnsec
			} else {
				pFsINodeMeta.Atime = input.Atime
				pFsINodeMeta.Atimensec = input.Atimensec
			}
		}

		if input.Valid&fsapi.FATTR_MTIME != 0 {
			if input.Valid&fsapi.FATTR_MTIME_NOW != 0 {
				pFsINodeMeta.Mtime = nowt
				pFsINodeMeta.Mtimensec = nowtnsec
			} else {
				pFsINodeMeta.Mtime = input.Mtime
				pFsINodeMeta.Mtimensec = input.Mtimensec
			}
		}
	}

	pFsINodeMeta.Ctime = nowt
	pFsINodeMeta.Ctimensec = nowtnsec

	return fsapi.OK
}

func (p *PosixFs) SetFsAttrOutByFsINode(fsAttrOut *fsapi.AttrOut,
	pFsINodeMeta *solofstypes.FsINodeMeta) {

	fsAttrOut.AttrValid = p.FsINodeDriver.EntryAttrValid
	fsAttrOut.AttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsAttrOut.Attr, pFsINodeMeta)
}

func (p *PosixFs) GetAttr(input *fsapi.GetAttrIn, out *fsapi.AttrOut) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapi.OK
}

func (p *PosixFs) SetAttr(input *fsapi.SetAttrIn, out *fsapi.AttrOut) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		code        fsapi.Status
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	code = p.SetFsINodeByFsAttr(&fsINodeMeta, input)
	if code != fsapi.OK {
		return code
	}

	err = p.FsINodeDriver.UpdateFsINode(&fsINodeMeta)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapi.OK
}
