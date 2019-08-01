package metastg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		soloOSEnv   soloosbase.SoloOSEnv
		peerPool    offheap.LKVTableWithBytes64
		metaStg     MetaStg
		netINode    sdfsapitypes.NetINode
		netBlock    sdfsapitypes.NetBlock
		netINodeID0 sdfsapitypes.NetINodeID
		netINodeID1 sdfsapitypes.NetINodeID
		netINodeID2 sdfsapitypes.NetINodeID
		peerID0     snettypes.PeerID
		peerID1     snettypes.PeerID
		err         error
	)
	util.AssertErrIsNil(soloOSEnv.InitWithSNet(""))

	util.AssertErrIsNil(metaStg.Init(&soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	sdfsapitypes.InitTmpNetINodeID(&netINodeID0)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID1)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID2)

	err = soloOSEnv.OffheapDriver.InitLKVTableWithBytes64(&peerPool, "TestMetaStgNet",
		int(snettypes.PeerStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	util.AssertErrIsNil(err)

	netINode.ID = netINodeID0
	netBlock.NetINodeID = netINode.ID

	snettypes.InitTmpPeerID(&peerID0)
	snettypes.InitTmpPeerID(&peerID1)

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
		assert.Equal(t, metaStg.FetchNetBlockFromDB(&netINode, 1, &netBlock, &backendPeerIDArrStr), sdfsapitypes.ErrObjectNotExists)
	}
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}

	util.AssertErrIsNil(metaStg.Close())
}
