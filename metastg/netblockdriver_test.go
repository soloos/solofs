package metastg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
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
		netINode    solofstypes.NetINode
		netBlock    solofstypes.NetBlock
		netINodeID0 solofstypes.NetINodeID
		netINodeID1 solofstypes.NetINodeID
		netINodeID2 solofstypes.NetINodeID
		peerID0     snet.PeerID
		peerID1     snet.PeerID
		err         error
	)
	util.AssertErrIsNil(soloosEnv.InitWithSNet(""))

	util.AssertErrIsNil(metaStg.Init(&soloosEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	solofstypes.InitTmpNetINodeID(&netINodeID0)
	solofstypes.InitTmpNetINodeID(&netINodeID1)
	solofstypes.InitTmpNetINodeID(&netINodeID2)

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
		assert.Equal(t, metaStg.FetchNetBlockFromDB(&netINode, 1, &netBlock, &backendPeerIDArrStr), solofstypes.ErrObjectNotExists)
	}
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}

	util.AssertErrIsNil(metaStg.Close())
}
