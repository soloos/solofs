package metastg

import (
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgINode(t *testing.T) {
	var (
		metastg MetaStg
		inode   types.INode
	)

	assert.NoError(t, metastg.Init(TestMetaStgDBDriver, TestMetaStgDBConnect))
	inode.ID = types.INodeID{1, 2, 4, 51, 25}
	assert.NoError(t, metastg.StoreINode(&inode))
	assert.NoError(t, metastg.StoreINode(&inode))

	{
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}
	{
		inode.ID = types.INodeID{1, 2, 4, 51}
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, false)
		assert.NoError(t, err)
	}
	{
		inode.ID = types.INodeID{1, 2, 4, 51, 25}
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
