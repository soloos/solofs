package metastg

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		soloOSEnv   soloosbase.SoloOSEnv
		metaStg     MetaStg
		netINode    types.NetINode
		netINodeID0 types.NetINodeID
		netINodeID1 types.NetINodeID
	)
	util.AssertErrIsNil(soloOSEnv.Init())

	assert.NoError(t, metaStg.Init(&soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	sdfsapitypes.InitTmpNetINodeID(&netINodeID0)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID1)

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
		assert.Equal(t, err, types.ErrObjectNotExists)
	}
	{
		netINode.ID = netINodeID0
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metaStg.Close())
}
