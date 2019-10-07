package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
)

type PosixFs struct {
	*soloosbase.SoloosEnv
	NameSpaceID solofstypes.NameSpaceID
	MemStg      *MemStg

	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
	FdTable       FdTable
	FsMutexDriver FsMutexDriver
}

var _ = fsapi.PosixFs(&PosixFs{})

func (p *PosixFs) Init(
	soloosEnv *soloosbase.SoloosEnv,
	nsID solofstypes.NameSpaceID,
	memStg *MemStg,
	// FsINodeDriver
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	getNetINode solofstypes.GetNetINode,
	mustGetNetINode solofstypes.MustGetNetINode,
	releaseNetINode solofstypes.ReleaseNetINode,
	allocFsINodeID solofstypes.AllocFsINodeID,
	deleteFsINodeByIDInDB solofstypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB solofstypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB solofstypes.UpdateFsINodeInDB,
	insertFsINodeInDB solofstypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB solofstypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB solofstypes.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB solofstypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofstypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofstypes.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.NameSpaceID = nsID
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
	)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(
		p,
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	err = p.FdTable.Init()
	if err != nil {
		return err
	}

	err = p.FsMutexDriver.Init(p.SoloosEnv, p)
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

func (p *PosixFs) Close() error {
	return nil
}
