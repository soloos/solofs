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
		metaStg       MetaStg
		netINode      types.NetINode
		netBlock      types.NetBlock
		id0           types.NetINodeID
		id1           types.NetINodeID
		id2           types.NetINodeID
	)

	util.AssertErrIsNil(metaStg.Init(offheapDriver, TestMetaStgDBDriver, TestMetaStgDBConnect, nil, nil))
	util.InitUUID64(&id0)
	util.InitUUID64(&id1)
	util.InitUUID64(&id2)

	util.AssertErrIsNil(offheap.DefaultOffheapDriver.InitRawObjectPool(&peerPool, int(snettypes.PeerStructSize), -1, nil, nil))

	netINode.ID = id0
	netBlock.NetINodeID = netINode.ID

	uPeer0 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer0.Ptr().PeerID)
	uPeer1 := snettypes.PeerUintptr(peerPool.AllocRawObject())
	util.InitUUID64(&uPeer1.Ptr().PeerID)
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
