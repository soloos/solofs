package metastg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetINode(t *testing.T) {
	var (
		soloOSEnv   soloosbase.SoloOSEnv
		metaStg     MetaStg
		netINode    solofsapitypes.NetINode
		netINodeID0 solofsapitypes.NetINodeID
		netINodeID1 solofsapitypes.NetINodeID
	)
	util.AssertErrIsNil(soloOSEnv.InitWithSNet(""))

	assert.NoError(t, metaStg.Init(&soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	solofsapitypes.InitTmpNetINodeID(&netINodeID0)
	solofsapitypes.InitTmpNetINodeID(&netINodeID1)

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
		assert.Equal(t, err, solofsapitypes.ErrObjectNotExists)
	}
	{
		netINode.ID = netINodeID0
		err := metaStg.FetchNetINodeFromDB(&netINode)
		assert.NoError(t, err)
	}

	assert.NoError(t, metaStg.Close())
}
