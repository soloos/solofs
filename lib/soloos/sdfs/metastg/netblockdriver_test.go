package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaStgNetBlock(t *testing.T) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		peerPool      offheap.RawObjectPool
		metastg       MetaStg
		netINode      types.NetINode
		netBlock      types.NetBlock
		id0           types.NetINodeID
		id1           types.NetINodeID
		id2           types.NetINodeID
	)

	assert.NoError(t, metastg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)
	util.InitUUID64(&id2)

	assert.NoError(t, offheap.DefaultOffheapDriver.InitRawObjectPool(&peerPool, int(snettypes.PeerStructSize), -1, nil, nil))

	netINode.ID = id0
	netBlock.ID = id1

	uPeer0 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer0.Ptr().PeerID)
	uPeer1 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer1.Ptr().PeerID)
	netBlock.DataNodes.Append(uPeer0)
	netBlock.DataNodes.Append(uPeer1)

	assert.NoError(t, metastg.StoreNetBlockInDB(&netINode, &netBlock))
	assert.NoError(t, metastg.StoreNetBlockInDB(&netINode, &netBlock))

	var backendPeerIDArrStr string
	{
		assert.NoError(t, metastg.FetchNetBlockFromDB(&netBlock, &backendPeerIDArrStr))
	}
	{
		netBlock.ID = id2
		assert.Equal(t, metastg.FetchNetBlockFromDB(&netBlock, &backendPeerIDArrStr), types.ErrObjectNotExists)
	}
	{
		netBlock.ID = id1
		assert.NoError(t, metastg.FetchNetBlockFromDB(&netBlock, &backendPeerIDArrStr))
	}

	assert.NoError(t, metastg.Close())
}
