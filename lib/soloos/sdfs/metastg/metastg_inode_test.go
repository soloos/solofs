package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgINode(t *testing.T) {
	var (
		metastg MetaStg
		inode   types.INode
		id0     types.INodeID
		id1     types.INodeID
	)

	assert.NoError(t, metastg.Init(TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)

	inode.ID = id0
	assert.NoError(t, metastg.StoreINode(&inode))
	assert.NoError(t, metastg.StoreINode(&inode))

	{
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}
	{
		inode.ID = id1
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, false)
		assert.NoError(t, err)
	}
	{
		inode.ID = id0
		exisits, err := metastg.FetchINode(&inode)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
