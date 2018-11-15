package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUintptrArray8(t *testing.T) {
	var arr UintptrArray8
	arr.Append(12)
	assert.Equal(t, arr.Len, 1)
	assert.Equal(t, arr.Arr[0], uintptr(12))
}
