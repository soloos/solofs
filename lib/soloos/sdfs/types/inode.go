package types

import (
	"sync"
	"unsafe"
)

const (
	INodeIDSize     = int(unsafe.Sizeof([64]byte{}))
	INodeStructSize = unsafe.Sizeof(INode{})
)

type INodeID = [INodeIDSize]byte
type INodeUintptr uintptr

func (u INodeUintptr) Ptr() *INode { return (*INode)(unsafe.Pointer(u)) }

type INode struct {
	ID           INodeID      `db:"inode_id"`
	Size         uint64       `db:"inode_size"`
	NetBlockSize int          `db:"netblock_size"`
	MemBlockSize int          `db:"memblock_size"`
	WriteRWMutex sync.RWMutex `db:"-"`
}

func (p *INode) IDStr() string { return string(p.ID[:]) }
