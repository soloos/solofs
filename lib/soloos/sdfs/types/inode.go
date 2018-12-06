package types

import (
	"sync"
	"unsafe"
)

const (
	INodeIDBytesNum = 64
	INodeIDSize     = int(unsafe.Sizeof([INodeIDBytesNum]byte{}))
	INodeStructSize = unsafe.Sizeof(INode{})
)

type INodeID = [INodeIDBytesNum]byte
type INodeUintptr uintptr

func (u INodeUintptr) Ptr() *INode { return (*INode)(unsafe.Pointer(u)) }

type INode struct {
	ID            INodeID      `db:"inode_id"`
	Size          uint64       `db:"inode_size"`
	NetBlockCap   int          `db:"netblock_cap"`
	MemBlockCap   int          `db:"memblock_cap"`
	AccessRWMutex sync.RWMutex `db:"-"`
}

func (p *INode) IDStr() string { return string(p.ID[:]) }
