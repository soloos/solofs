package metastg

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solodb/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		soloosEnv   soloosbase.SoloosEnv
		peerPool    offheap.LKVTableWithBytes64
		metaStg     MetaStg
		netINode    solofsapitypes.NetINode
		netBlock    solofsapitypes.NetBlock
		netINodeID0 solofsapitypes.NetINodeID
		netINodeID1 solofsapitypes.NetINodeID
		netINodeID2 solofsapitypes.NetINodeID
		peerID0     snet.PeerID
		peerID1     snet.PeerID
		err         error
	)
	util.AssertErrIsNil(soloosEnv.InitWithSNet(""))

	util.AssertErrIsNil(metaStg.Init(&soloosEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	solofsapitypes.InitTmpNetINodeID(&netINodeID0)
	solofsapitypes.InitTmpNetINodeID(&netINodeID1)
	solofsapitypes.InitTmpNetINodeID(&netINodeID2)

	err = soloosEnv.OffheapDriver.InitLKVTableWithBytes64(&peerPool, "TestMetaStgNet",
		int(snet.PeerStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	util.AssertErrIsNil(err)

	netINode.ID = netINodeID0
	netBlock.NetINodeID = netINode.ID

	snet.InitTmpPeerID(&peerID0)
	snet.InitTmpPeerID(&peerID1)

	netBlock.StorDataBackends.Append(peerID0)
	netBlock.StorDataBackends.Append(peerID1)
	netBlock.IndexInNetINode = 0

	util.AssertErrIsNil(metaStg.StoreNetBlockInDB(&netINode, &netBlock))
	util.AssertErrIsNil(metaStg.StoreNetBlockInDB(&netINode, &netBlock))

	var backendPeerIDArrStr string
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}
	{
		assert.Equal(t, metaStg.FetchNetBlockFromDB(&netINode, 1, &netBlock, &backendPeerIDArrStr), solofsapitypes.ErrObjectNotExists)
	}
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}

	util.AssertErrIsNil(metaStg.Close())
}
