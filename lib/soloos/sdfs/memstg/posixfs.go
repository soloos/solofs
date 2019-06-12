package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"time"
)

type PosixFS struct {
	*soloosbase.SoloOSEnv
	NameSpaceID   sdfsapitypes.NameSpaceID
	MemStg        *MemStg
	FsINodeDriver FsINodeDriver
	FdTable       FdTable

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32
}

var _ = fsapi.PosixFS(&PosixFS{})

func (p *PosixFS) Init(
	soloOSEnv *soloosbase.SoloOSEnv,
	nameSpaceID sdfsapitypes.NameSpaceID,
	memStg *MemStg,
	// FsINodeDriver
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID sdfsapitypes.AllocFsINodeID,
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	releaseNetINode sdfsapitypes.ReleaseNetINode,
	deleteFsINodeByIDInDB sdfsapitypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB sdfsapitypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB sdfsapitypes.UpdateFsINodeInDB,
	insertFsINodeInDB sdfsapitypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB sdfsapitypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB sdfsapitypes.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB sdfsapitypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB sdfsapitypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB sdfsapitypes.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.NameSpaceID = nameSpaceID
	p.MemStg = memStg

	err = p.FsINodeDriver.Init(p.SoloOSEnv,
		p,
		defaultNetBlockCap,
		defaultMemBlockCap,
		allocFsINodeID,
		getNetINode,
		mustGetNetINode,
		releaseNetINode,
		deleteFsINodeByIDInDB,
		listFsINodeByParentIDFromDB,
		updateFsINodeInDB,
		insertFsINodeInDB,
		fetchFsINodeByIDFromDB,
		fetchFsINodeByNameFromDB,
		// FIXAttrDriver
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	p.refreshEntryTtl()

	err = p.FdTable.Init()
	if err != nil {
		return err
	}

	return nil
}

// This is called on processing the first request. The
// filesystem implementation can use the server argument to
// talk back to the kernel (through notify methods).
func (p *PosixFS) FsInit() {
}

func (p *PosixFS) String() string {
	return "sdfs"
}

// If called, provide debug output through the log package.
func (p *PosixFS) SetDebug(debug bool) {
}

func (p *PosixFS) refreshEntryTtl() {
	p.EntryTtl = p.FsINodeDriver.EntryTtl
	p.EntryAttrValid = p.FsINodeDriver.EntryAttrValid
	p.EntryAttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
}
