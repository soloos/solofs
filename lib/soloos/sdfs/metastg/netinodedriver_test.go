package metastg

import (
	"soloos/sdfs/types"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		metaStg       MetaStg
		netINode      types.NetINode
		id0           types.NetINodeID
		id1           types.NetINodeID
	)

	assert.NoError(t, metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)

	netINode.ID = id0
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))

	{
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, nil)
	}
	{
		netINode.ID = id1
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, types.ErrObjectNotExists)
	}
	{
		netINode.ID = id0
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metaStg.Close())
}
