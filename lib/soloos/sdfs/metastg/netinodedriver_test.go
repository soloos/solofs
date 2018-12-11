package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		metastg       MetaStg
		netINode         types.NetINode
		id0           types.NetINodeID
		id1           types.NetINodeID
	)

	assert.NoError(t, metastg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)

	netINode.ID = id0
	assert.NoError(t, metastg.StoreNetINodeInDB(&netINode))
	assert.NoError(t, metastg.StoreNetINodeInDB(&netINode))

	{
		err := metastg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, nil)
	}
	{
		netINode.ID = id1
		err := metastg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, types.ErrObjectNotExists)
	}
	{
		netINode.ID = id0
		err := metastg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
