package metastg

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		metaStg       MetaStg
		netINode      types.NetINode
		peerID0       types.NetINodeID
		peerID1       types.NetINodeID
	)

	assert.NoError(t, metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	snettypes.InitTmpPeerID(&peerID0)
	snettypes.InitTmpPeerID(&peerID1)

	netINode.ID = peerID0
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))

	{
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, nil)
	}
	{
		netINode.ID = peerID1
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, types.ErrObjectNotExists)
	}
	{
		netINode.ID = peerID0
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metaStg.Close())
}
