package types

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestEncodeINodeBlockID(t *testing.T) {
	var inodeBlockID INodeBlockID
	inodeID := INodeID{1, 2, 3}
	blockIndex := 22
	EncodeINodeBlockID(&inodeBlockID, inodeID, blockIndex)
	assert.Equal(t, uint8(1), inodeBlockID[0])
	assert.Equal(t, uint8(2), inodeBlockID[1])
	assert.Equal(t, uint8(3), inodeBlockID[2])
	assert.Equal(t, uint8(22), inodeBlockID[INodeIDSize])
}

func TestEncodePtrBindIndex(t *testing.T) {
	var (
		u     uintptr = 0x12
		index int     = 3
		id    PtrBindIndex
	)
	EncodePtrBindIndex(&id, u, index)
	assert.Equal(t, uintptr(0x12), *((*uintptr)(unsafe.Pointer(&id))))
	assert.Equal(t, uint8(3), id[UintptrSize])
}

func BenchmarkEncodeINodeBlockID(b *testing.B) {
	var inodeBlockID INodeBlockID
	inodeID := INodeID{1, 2, 3}
	blockIndex := 22
	for n := 0; n < b.N; n++ {
		EncodeINodeBlockID(&inodeBlockID, inodeID, blockIndex)
	}
}

func BenchmarkEncodeBindIndex(b *testing.B) {
	var (
		u     uintptr = 0x12
		index int     = 3
		id    PtrBindIndex
	)
	for n := 0; n < b.N; n++ {
		EncodePtrBindIndex(&id, u, index)
	}
}
