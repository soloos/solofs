package metastg

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		soloOSEnv   soloosbase.SoloOSEnv
		peerPool    offheap.LKVTableWithBytes64
		uObject     offheap.LKVTableObjectUPtrWithBytes64
		metaStg     MetaStg
		netINode    types.NetINode
		netBlock    types.NetBlock
		netINodeID0 types.NetINodeID
		netINodeID1 types.NetINodeID
		netINodeID2 types.NetINodeID
		peerID      snettypes.PeerID
		err         error
	)
	util.AssertErrIsNil(soloOSEnv.Init())

	util.AssertErrIsNil(metaStg.Init(&soloOSEnv, TestMetaStgDBDriver, TestMetaStgDBConnect))
	sdfsapitypes.InitTmpNetINodeID(&netINodeID0)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID1)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID2)

	err = soloOSEnv.OffheapDriver.InitLKVTableWithBytes64(&peerPool, "TestMetaStgNet",
		int(snettypes.PeerStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	util.AssertErrIsNil(err)

	netINode.ID = netINodeID0
	netBlock.NetINodeID = netINode.ID

	snettypes.InitTmpPeerID(&peerID)
	uObject, _ = peerPool.MustGetObjectWithAcquire(peerID)
	uPeer0 := snettypes.PeerUintptr(uObject)

	snettypes.InitTmpPeerID(&peerID)
	uObject, _ = peerPool.MustGetObjectWithAcquire(peerID)
	uPeer1 := snettypes.PeerUintptr(uObject)

	netBlock.StorDataBackends.Append(uPeer0)
	netBlock.StorDataBackends.Append(uPeer1)
	netBlock.IndexInNetINode = 0

	util.AssertErrIsNil(metaStg.StoreNetBlockInDB(&netINode, &netBlock))
	util.AssertErrIsNil(metaStg.StoreNetBlockInDB(&netINode, &netBlock))

	var backendPeerIDArrStr string
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}
	{
		assert.Equal(t, metaStg.FetchNetBlockFromDB(&netINode, 1, &netBlock, &backendPeerIDArrStr), types.ErrObjectNotExists)
	}
	{
		util.AssertErrIsNil(metaStg.FetchNetBlockFromDB(&netINode, 0, &netBlock, &backendPeerIDArrStr))
	}

	util.AssertErrIsNil(metaStg.Close())
}
