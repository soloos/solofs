package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgINode(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		metastg       MetaStg
		inode         types.INode
		id0           types.INodeID
		id1           types.INodeID
	)

	assert.NoError(t, metastg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)

	inode.ID = id0
	assert.NoError(t, metastg.StoreINodeInDB(&inode))
	assert.NoError(t, metastg.StoreINodeInDB(&inode))

	{
		err := metastg.FetchINodeFromDB(&inode)
		assert.Equal(t, err, nil)
	}
	{
		inode.ID = id1
		err := metastg.FetchINodeFromDB(&inode)
		assert.Equal(t, err, types.ErrObjectNotExists)
	}
	{
		inode.ID = id0
		err := metastg.FetchINodeFromDB(&inode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metastg.Close())
}
