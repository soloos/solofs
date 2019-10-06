package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"time"
)

type PosixFs struct {
	*soloosbase.SoloosEnv
	NameSpaceID solofsapitypes.NameSpaceID
	MemStg      *MemStg

	FsINodeDriver
	FdTable

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32
}

var _ = fsapi.PosixFs(&PosixFs{})

func (p *PosixFs) Init(
	soloosEnv *soloosbase.SoloosEnv,
	nameSpaceID solofsapitypes.NameSpaceID,
	memStg *MemStg,
	// FsINodeDriver
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID solofsapitypes.AllocFsINodeID,
	getNetINode solofsapitypes.GetNetINode,
	mustGetNetINode solofsapitypes.MustGetNetINode,
	releaseNetINode solofsapitypes.ReleaseNetINode,
	deleteFsINodeByIDInDB solofsapitypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB solofsapitypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB solofsapitypes.UpdateFsINodeInDB,
	insertFsINodeInDB solofsapitypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB solofsapitypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB solofsapitypes.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB solofsapitypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofsapitypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofsapitypes.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.NameSpaceID = nameSpaceID
	p.MemStg = memStg

	err = p.FsINodeDriver.Init(p.SoloosEnv,
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
func (p *PosixFs) FsInit() {
}

func (p *PosixFs) String() string {
	return "solofs"
}

// If called, provide debug output through the log package.
func (p *PosixFs) SetDebug(debug bool) {
}

func (p *PosixFs) refreshEntryTtl() {
	p.EntryTtl = p.FsINodeDriver.EntryTtl
	p.EntryAttrValid = p.FsINodeDriver.EntryAttrValid
	p.EntryAttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
}
