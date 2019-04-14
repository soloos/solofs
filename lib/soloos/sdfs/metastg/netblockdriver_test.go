package metastg

import (
	snettypes "soloos/common/snet/types"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		peerPool      offheap.LKVTableWithBytes64
		uObject       uintptr
		metaStg       MetaStg
		netINode      types.NetINode
		netBlock      types.NetBlock
		id0           types.NetINodeID
		id1           types.NetINodeID
		id2           types.NetINodeID
		peerID        snettypes.PeerID
		err           error
	)

	util.AssertErrIsNil(metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)
	util.InitUUID64(&id2)

	err = offheap.DefaultOffheapDriver.InitLKVTableWithBytes64(&peerPool, "TestMetaStgNet",
		int(snettypes.PeerStructSize), -1, types.DefaultKVTableSharedCount, nil, nil)
	util.AssertErrIsNil(err)

	netINode.ID = id0
	netBlock.NetINodeID = netINode.ID

	util.InitUUID64(&peerID)
	uObject, _ = peerPool.MustGetObjectWithAcquire(peerID)
	uPeer0 := snettypes.PeerUintptr(uObject)

	util.InitUUID64(&peerID)
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
