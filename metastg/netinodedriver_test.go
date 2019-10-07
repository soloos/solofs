package metastg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		soloosEnv   soloosbase.SoloosEnv
		metaStg     MetaStg
		netINode    solofstypes.NetINode
		netINodeID0 solofstypes.NetINodeID
		netINodeID1 solofstypes.NetINodeID
	)
	util.AssertErrIsNil(soloosEnv.InitWithSNet(""))

	assert.NoError(t, metaStg.Init(&soloosEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	solofstypes.InitTmpNetINodeID(&netINodeID0)
	solofstypes.InitTmpNetINodeID(&netINodeID1)

	netINode.ID = netINodeID0
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))
	assert.NoError(t, metaStg.StoreNetINodeInDB(&netINode))

	{
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, nil)
	}
	{
		netINode.ID = netINodeID1
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.Equal(t, err, solofstypes.ErrObjectNotExists)
	}
	{
		netINode.ID = netINodeID0
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metaStg.Close())
}
