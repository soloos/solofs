package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		metastg  MetaStg
		inode    types.INode
		netBlock types.NetBlock
		id0      types.INodeID
		id1      types.INodeID
		id2      types.INodeID
	)

	assert.NoError(t, metastg.Init(TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)
	util.InitUUID64(&id2)

	inode.ID = id0
	netBlock.ID = id1
	assert.NoError(t, metastg.StoreNetBlock(&inode, &netBlock, true))
	assert.NoError(t, metastg.StoreNetBlock(&inode, &netBlock, true))

	{
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}
	{
		netBlock.ID = id2
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, false)
		assert.NoError(t, err)
	}
	{
		netBlock.ID = id1
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
