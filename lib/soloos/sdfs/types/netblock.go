package types

import (
	snettypes "soloos/snet/types"
	"sync"
	"unsafe"
)

const (
	NetBlockStructSize            = unsafe.Sizeof(NetBlock{})
	MaxDataNodesSizeStoreNetBlock = 8
)

type NetBlockUintptr uintptr

func (u NetBlockUintptr) Ptr() *NetBlock { return (*NetBlock)(unsafe.Pointer(u)) }

type NetBlock struct {
	NetINodeID       NetINodeID                  `db:"netinode_id"`
	IndexInNetINode  int                         `db:"index_in_netinode"`
	Len              int                         `db:"netblock_len"`
	Cap              int                         `db:"netblock_cap"`
	DataNodes        snettypes.PeerUintptrArray8 `db:"-"`
	MetaDataMutex    sync.Mutex                  `db:"-"`
	IsMetaDataInited bool                        `db:"-"`
}

func (p *NetBlock) NetINodeIDStr() string { return string(p.NetINodeID[:]) }
