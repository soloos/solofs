package metastg

import (
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		metastg  MetaStg
		inode    types.INode
		netBlock types.NetBlock
	)

	assert.NoError(t, metastg.Init(TestMetaStgDBDriver, TestMetaStgDBConnect))
	inode.ID = types.INodeID{1, 2, 4, 51, 25}
	netBlock.ID = types.NetBlockID{1, 2, 4, 51, 25}
	assert.NoError(t, metastg.StoreNetBlock(&inode, &netBlock, true))
	assert.NoError(t, metastg.StoreNetBlock(&inode, &netBlock, true))

	{
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}
	{
		netBlock.ID = types.NetBlockID{1, 2, 4, 51}
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, false)
		assert.NoError(t, err)
	}
	{
		netBlock.ID = types.NetBlockID{1, 2, 4, 51, 25}
		exisits, err := metastg.FetchNetBlock(&netBlock)
		assert.Equal(t, exisits, true)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
