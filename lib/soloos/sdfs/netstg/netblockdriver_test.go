package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		offheapDriver         = &offheap.DefaultOffheapDriver
		netBlockDriverOptions = NetBlockDriverOptions{
			NetBlockPoolOptions{
				int32(-1),
			},
		}
		netBlockDriver   NetBlockDriver
		snetClientDriver snet.ClientDriver
	)

	assert.NoError(t, snetClientDriver.Init(offheapDriver))

	var peer snettypes.Peer
	peer.ID = snettypes.PeerID{}
	copy(peer.Address[:], []byte(MockSdfsdAddress))
	copy(peer.ServiceProtocol[:], []byte("srpc"))
	uPeer := snetClientDriver.RegisterPeer(peer)

	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, &snetClientDriver))

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	var (
		inode  types.INode
		uINode types.INodeUintptr = types.INodeUintptr((unsafe.Pointer(&inode)))
	)
	uINode.Ptr().NetBlockSize = 1024
	uINode.Ptr().MemBlockSize = 1024

	uNetBlock, _ := netBlockDriver.MustGetBlock(uINode, 10)
	uNetBlock.Ptr().DataNodes.Append(uPeer)
	assert.NoError(t, netBlockDriver.WriteBytesAt(uINode, uNetBlock, data, 10))
}
