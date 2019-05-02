package types

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
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
	offheap.LKVTableObjectWithBytes68 `db:"-"`

	NetINodeID      NetINodeID `db:"netinode_id"`
	IndexInNetINode int32      `db:"index_in_netinode"`
	Len             int        `db:"netblock_len"`
	Cap             int        `db:"netblock_cap"`

	StorDataBackends    snettypes.PeerGroup `db:"-"`
	IsDBMetaDataInited  MetaDataState               `db:"-"`
	DBMetaDataInitMutex sync.Mutex                  `db:"-"`

	SyncDataBackends                    snettypes.PeerGroup `db:"-"`
	SyncDataPrimaryBackendTransferCount int                         `db:"-"`
	IsSyncDataBackendsInited            MetaDataState               `db:"-"`
	LocalDataBackend                    snettypes.PeerUintptr       `db:"-"`
	IsLocalDataBackendInited            MetaDataState               `db:"-"`
	MemMetaDataInitMutex                sync.Mutex                  `db:"-"`
}

func (p *NetBlock) NetINodeIDStr() string { return string(p.NetINodeID[:]) }

func (p *NetBlock) Reset() {
	p.IsDBMetaDataInited.Reset()
	p.IsSyncDataBackendsInited.Reset()
	p.IsLocalDataBackendInited.Reset()
}
