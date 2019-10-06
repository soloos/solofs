package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFs) SetFsEntryOutByFsINode(fsEntryOut *fsapitypes.EntryOut,
	pFsINodeMeta *solofsapitypes.FsINodeMeta) {

	fsEntryOut.NodeId = pFsINodeMeta.Ino
	fsEntryOut.EntryValid = p.EntryAttrValid
	fsEntryOut.EntryValidNsec = p.EntryAttrValidNsec
	fsEntryOut.AttrValid = p.EntryAttrValid
	fsEntryOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsEntryOut.Attr, pFsINodeMeta)
}

func (p *PosixFs) SetFsINodeByFsAttr(pFsINodeMeta *solofsapitypes.FsINodeMeta,
	input *fsapitypes.SetAttrIn) fsapitypes.Status {

	if input.Valid&fsapitypes.FATTR_MODE != 0 {
		// pFsINodeMeta.Mode = uint32(0777)&input.Mode | uint32(solofsapitypes.FsINodeTypeToFsType(pFsINodeMeta.Type))
		pFsINodeMeta.Mode = input.Mode
	}

	if input.Valid&fsapitypes.FATTR_UID != 0 {
		pFsINodeMeta.Uid = input.Uid
	}
	if input.Valid&fsapitypes.FATTR_GID != 0 {
		pFsINodeMeta.Gid = input.Gid
	}

	if input.Valid&fsapitypes.FATTR_SIZE != 0 {
		if pFsINodeMeta.Type == solofstypes.FSINODE_TYPE_DIR {
			return fsapitypes.EISDIR
		}
		var err = p.TruncateINode(pFsINodeMeta, input.Size)
		if err != nil {
			return fsapitypes.EIO
		}
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := solofsapitypes.DirTreeTime(now.Unix())
	nowtnsec := solofsapitypes.DirTreeTimeNsec(now.UnixNano())

	if input.Valid&(fsapitypes.FATTR_ATIME|fsapitypes.FATTR_MTIME|fsapitypes.FATTR_ATIME_NOW|fsapitypes.FATTR_MTIME_NOW) != 0 {
		if input.Valid&fsapitypes.FATTR_ATIME != 0 {
			if input.Valid&fsapitypes.FATTR_ATIME_NOW != 0 {
				pFsINodeMeta.Atime = nowt
				pFsINodeMeta.Atimensec = nowtnsec
			} else {
				pFsINodeMeta.Atime = input.Atime
				pFsINodeMeta.Atimensec = input.Atimensec
			}
		}

		if input.Valid&fsapitypes.FATTR_MTIME != 0 {
			if input.Valid&fsapitypes.FATTR_MTIME_NOW != 0 {
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

	return fsapitypes.OK
}

func (p *PosixFs) SetFsAttrOutByFsINode(fsAttrOut *fsapitypes.AttrOut,
	pFsINodeMeta *solofsapitypes.FsINodeMeta) {

	fsAttrOut.AttrValid = p.EntryAttrValid
	fsAttrOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsAttrOut.Attr, pFsINodeMeta)
}

func (p *PosixFs) GetAttr(input *fsapitypes.GetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFs) SetAttr(input *fsapitypes.SetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		code        fsapitypes.Status
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	code = p.SetFsINodeByFsAttr(&fsINodeMeta, input)
	if code != fsapitypes.OK {
		return code
	}

	err = p.FsINodeDriver.UpdateFsINode(&fsINodeMeta)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}
