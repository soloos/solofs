package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"time"

	"github.com/hanwen/go-fuse/fuse"
)

type DirTreeStg struct {
	MemStg        *MemStg
	FsINodeDriver FsINodeDriver
	FdTable       FdTable

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32
}

var _ = fuse.RawFileSystem(&DirTreeStg{})

// This is called on processing the first request. The
// filesystem implementation can use the server argument to
// talk back to the kernel (through notify methods).
func (p *DirTreeStg) Init(server *fuse.Server) {
}

func (p *DirTreeStg) String() string {
	return types.FuseName
}

// If called, provide debug output through the log package.
func (p *DirTreeStg) SetDebug(debug bool) {
}

func (p *DirTreeStg) SdfsInit(
	memStg *MemStg,
	offheapDriver *offheap.OffheapDriver,
	// FsINodeDriver
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID api.AllocFsINodeID,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	getFsINodeByIDFromDB api.GetFsINodeByIDFromDB,
	getFsINodeByNameFromDB api.GetFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.MemStg = memStg

	err = p.FsINodeDriver.Init(
		offheapDriver,
		p,
		defaultNetBlockCap,
		defaultMemBlockCap,
		allocFsINodeID,
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire,
		deleteFsINodeByIDInDB,
		listFsINodeByParentIDFromDB,
		updateFsINodeInDB,
		insertFsINodeInDB,
		getFsINodeByIDFromDB,
		getFsINodeByNameFromDB,
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

func (p *DirTreeStg) refreshEntryTtl() {
	p.EntryTtl = p.FsINodeDriver.EntryTtl
	p.EntryAttrValid = p.FsINodeDriver.EntryAttrValid
	p.EntryAttrValidNsec = p.FsINodeDriver.EntryAttrValidNsec
}