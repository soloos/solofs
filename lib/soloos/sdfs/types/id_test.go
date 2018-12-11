package types

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestEncodeNetINodeBlockID(t *testing.T) {
	var netINodeBlockID NetINodeBlockID
	netINodeID := NetINodeID{1, 2, 3}
	blockIndex := 22
	EncodeNetINodeBlockID(&netINodeBlockID, netINodeID, blockIndex)
	assert.Equal(t, uint8(1), netINodeBlockID[0])
	assert.Equal(t, uint8(2), netINodeBlockID[1])
	assert.Equal(t, uint8(3), netINodeBlockID[2])
	assert.Equal(t, uint8(22), netINodeBlockID[NetINodeIDSize])
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

func BenchmarkEncodeNetINodeBlockID(b *testing.B) {
	var netINodeBlockID NetINodeBlockID
	netINodeID := NetINodeID{1, 2, 3}
	blockIndex := 22
	for n := 0; n < b.N; n++ {
		EncodeNetINodeBlockID(&netINodeBlockID, netINodeID, blockIndex)
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
