package memstg

import (
	"soloos/common/fsapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
	"time"
)

type DirTreeStg struct {
	*soloosbase.SoloOSEnv
	MemStg        *MemStg
	FsINodeDriver FsINodeDriver
	FdTable       FdTable

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32
}

var _ = fsapi.PosixFS(&DirTreeStg{})

func (p *DirTreeStg) Init(
	soloOSEnv *soloosbase.SoloOSEnv,
	memStg *MemStg,
	// FsINodeDriver
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID api.AllocFsINodeID,
	getNetINode api.GetNetINode,
	mustGetNetINode api.MustGetNetINode,
	releaseNetINode api.ReleaseNetINode,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB api.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB api.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
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
func (p *DirTreeStg) FsInit() {
}

func (p *DirTreeStg) String() string {
	return "sdfs"
}

// If called, provide debug output through the log package.
func (p *DirTreeStg) SetDebug(debug bool) {
}

func (p *DirTreeStg) refreshEntryTtl() {
	p.EntryTtl = p.FsINodeDriver.EntryTtl
	p.EntryAttrValid = p.FsINodeDriver.EntryAttrValid
	p.EntryAttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
}
