package metastg

import "C"
import (
	"soloos/dbcli"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"sync"
	"time"
)

type DirTreeDriverHelper struct {
	DBConn                         *dbcli.Connection
	OffheapDriver                  *offheap.OffheapDriver
	GetNetINodeWithReadAcquire     api.GetNetINodeWithReadAcquire
	MustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire
}

type DirTreeDriver struct {
	sysFsINode [2]types.FsINode

	allocINodeIDDalta types.FsINodeID
	lastFsINodeIDInDB types.FsINodeID
	maxFsINodeID      types.FsINodeID
	helper            DirTreeDriverHelper

	fsINodesByPathRWMutex sync.RWMutex
	fsINodesByPath        map[string]types.FsINode
	fsINodesByIDRWMutex   sync.RWMutex
	fsINodesByID          map[types.FsINodeID]types.FsINode
	rootFsINode           types.FsINode

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32
}

func (p *DirTreeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbcli.Connection,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) error {
	var err error
	p.SetHelper(offheapDriver,
		dbConn,
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire,
	)

	err = p.PrepareSchema()
	if err != nil {
		return err
	}

	err = p.PrepareINodes()
	if err != nil {
		return err
	}

	p.EntryTtl = 100 * time.Millisecond
	splitDuration(p.EntryTtl, &p.EntryAttrValid, &p.EntryAttrValidNsec)

	return nil
}

func (p *DirTreeDriver) SetHelper(offheapDriver *offheap.OffheapDriver,
	dbConn *dbcli.Connection,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) {
	p.helper.OffheapDriver = offheapDriver
	p.helper.DBConn = dbConn
	p.helper.GetNetINodeWithReadAcquire = getNetINodeWithReadAcquire
	p.helper.MustGetNetINodeWithReadAcquire = mustGetNetINodeWithReadAcquire
}
