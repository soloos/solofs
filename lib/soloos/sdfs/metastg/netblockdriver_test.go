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
		inode         types.INode
		netBlock      types.NetBlock
		id0           types.INodeID
		id1           types.INodeID
		id2           types.INodeID
	)

	assert.NoError(t, metastg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)
	util.InitUUID64(&id2)

	assert.NoError(t, offheap.DefaultOffheapDriver.InitRawObjectPool(&peerPool, int(snettypes.PeerStructSize), -1, nil, nil))

	inode.ID = id0
	netBlock.ID = id1

	uPeer0 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer0.Ptr().PeerID)
	uPeer1 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer1.Ptr().PeerID)
	netBlock.DataNodes.Append(uPeer0)
	netBlock.DataNodes.Append(uPeer1)

	assert.NoError(t, metastg.StoreNetBlockInDB(&inode, &netBlock))
	assert.NoError(t, metastg.StoreNetBlockInDB(&inode, &netBlock))

	{
		assert.NoError(t, metastg.FetchNetBlockFromDB(&netBlock))
	}
	{
		netBlock.ID = id2
		assert.Equal(t, metastg.FetchNetBlockFromDB(&netBlock), types.ErrObjectNotExists)
	}
	{
		netBlock.ID = id1
		assert.NoError(t, metastg.FetchNetBlockFromDB(&netBlock))
	}

	assert.NoError(t, metastg.Close())
}
