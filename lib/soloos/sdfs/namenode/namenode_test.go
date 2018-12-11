package namenode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MakeNameNodeForTest(nameNode *NameNode, nameNodeSRPCServerAddr string) {
	var (
		offheapDriver *offheap.OffheapDriver = &offheap.DefaultOffheapDriver

		options = NameNodeOptions{
			SRPCServer: NameNodeSRPCServerOptions{
				Network:    "tcp",
				ListenAddr: nameNodeSRPCServerAddr,
			},
			MetaStgDBDriver:  metastg.TestMetaStgDBDriver,
			MetaStgDBConnect: metastg.TestMetaStgDBConnect,
		}
		err error
	)
	err = nameNode.Init(options, offheapDriver)
	util.AssertErrIsNil(err)
}

func TestNetBlockPrepareMetadata(t *testing.T) {
	var (
		nameNode               NameNode
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
	)
	MakeNameNodeForTest(&nameNode, nameNodeSRPCListenAddr)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	var (
		memBlockDriver   memstg.MemBlockDriver
		netINodeDriver   memstg.NetINodeDriver
		netBlockCap      int   = 128
		memBlockCap      int   = 64
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		err              error
	)
	memstg.InitDriversForTest(t,
		nameNodeSRPCListenAddr,
		&memBlockDriver, &netINodeDriver, memBlockCap, blockChunksLimit)

	uNetINode, err = netINodeDriver.InitNetINode(0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)
	assert.NoError(t, netINodeDriver.PRead(uNetINode, readData, readOff))

	assert.NoError(t, nameNode.Close())
}
