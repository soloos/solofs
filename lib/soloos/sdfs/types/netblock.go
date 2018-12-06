package types

import (
	snettypes "soloos/snet/types"
	"sync"
	"unsafe"
)

const (
	NetBlockIDSize                = int(unsafe.Sizeof([64]byte{}))
	NetBlockStructSize            = unsafe.Sizeof(NetBlock{})
	MaxDataNodesSizeStoreNetBlock = 8
)

type NetBlockID = [NetBlockIDSize]byte
type NetBlockUintptr uintptr

func (u NetBlockUintptr) Ptr() *NetBlock { return (*NetBlock)(unsafe.Pointer(u)) }

type NetBlock struct {
	ID               NetBlockID `db:"netblock_id"`
	IndexInInode     int        `db:"index_in_inode"`
	Len              int32      `db:"netblock_len"`
	Cap              int32      `db:"netblock_cap"`
	MetaDataMutex    sync.Mutex
	IsMetaDataInited bool                        `db:"-"`
	DataNodes        snettypes.PeerUintptrArray8 `db:"-"`
}

func (p *NetBlock) IDStr() string { return string(p.ID[:]) }
