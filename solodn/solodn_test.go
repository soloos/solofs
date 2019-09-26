package solodn

import (
	"fmt"
	"soloos/common/iron"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"soloos/solofs/solonn"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	go util.PProfServe("192.168.56.100:17221")
	var (
		solonnIns                 solonn.Solonn
		solonnSRPCPeerID          = snet.MakeSysPeerID("SolonnSRPCForTest")
		solonnWebPeerID           = snet.MakeSysPeerID("SolonnWebForTest")
		solonnSRPCListenAddr      = "127.0.0.1:10401"
		solonnWebListenAddr       = "127.0.0.1:10402"
		netDriverWebServer        iron.Server
		netDriverServerListenAddr = "127.0.0.1:10403"
		netDriverServerServeAddr  = "http://127.0.0.1:10403"
		metaStgForSolonn          metastg.MetaStg

		solodns               [6]Solodn
		solodnSRPCPeerIDs     [6]snettypes.PeerID
		solodnSRPCListenAddrs = []string{
			"127.0.0.1:10410",
			"127.0.0.1:10411",
			"127.0.0.1:10412",
			"127.0.0.1:10413",
			"127.0.0.1:10414",
			"127.0.0.1:10415",
		}
	)

	var (
		soloosEnvForClient      soloosbase.SoloosEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient memstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloosEnvForSolonn      soloosbase.SoloosEnv
		memBlockDriverForSolonn memstg.MemBlockDriver
		netBlockDriverForSolonn memstg.NetBlockDriver
		netINodeDriverForSolonn memstg.NetINodeDriver

		soloosEnvForSolodns      [6]soloosbase.SoloosEnv
		memBlockDriverForSolodns [6]memstg.MemBlockDriver
		netBlockDriverForSolodns [6]memstg.NetBlockDriver
		netINodeDriverForSolodns [6]memstg.NetINodeDriver

		netBlockCap int   = 32
		memBlockCap int   = 16
		blocksLimit int32 = 4
		uNetINode   solofsapitypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloosEnvForSolonn.InitWithSNet(""))
	{
		var webServerOptions = iron.Options{
			ListenStr: netDriverServerListenAddr,
			ServeStr:  netDriverServerServeAddr,
		}
		util.AssertErrIsNil(netDriverWebServer.Init(webServerOptions))
	}
	go func() {
		assert.NoError(t, soloosEnvForSolonn.SNetDriver.PrepareServer("",
			&netDriverWebServer,
			nil, nil))
		assert.NoError(t, soloosEnvForSolonn.SNetDriver.ServerServe())
	}()
	time.Sleep(100 * time.Millisecond)
	metastg.MakeMetaStgForTest(&soloosEnvForSolonn, &metaStgForSolonn)

	assert.NoError(t, soloosEnvForClient.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloosEnvForClient,
		solonnSRPCListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloosEnvForSolonn,
		solonnSRPCListenAddr,
		&memBlockDriverForSolonn, &netBlockDriverForSolonn, &netINodeDriverForSolonn, memBlockCap, blocksLimit)
	solonn.MakeSolonnForTest(&soloosEnvForSolonn, &solonnIns, &metaStgForSolonn,
		solonnSRPCPeerID, solonnSRPCListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForSolonn, &netBlockDriverForSolonn, &netINodeDriverForSolonn)

	go func() {
		assert.NoError(t, solonnIns.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < len(solodnSRPCListenAddrs); i++ {
		assert.NoError(t, soloosEnvForSolodns[i].InitWithSNet(netDriverServerServeAddr))
		solodnSRPCPeerIDs[i] = snet.MakeSysPeerID(fmt.Sprintf("SolodnForTest_%v", i))

		memstg.MemStgMakeDriversForTest(&soloosEnvForSolodns[i],
			solonnSRPCListenAddr,
			&memBlockDriverForSolodns[i],
			&netBlockDriverForSolodns[i],
			&netINodeDriverForSolodns[i],
			memBlockCap, blocksLimit)

		MakeSolodnForTest(&soloosEnvForSolodns[i],
			&solodns[i],
			solodnSRPCPeerIDs[i], solodnSRPCListenAddrs[i],
			solonnSRPCPeerID, solonnSRPCListenAddr,
			&memBlockDriverForSolodns[i],
			&netBlockDriverForSolodns[i],
			&netINodeDriverForSolodns[i])
		go func(localI int) {
			assert.NoError(t, solodns[localI].Serve())
		}(i)
	}
	time.Sleep(time.Millisecond * 300)

	var (
		netINodeID solofsapitypes.NetINodeID
	)
	solofsapitypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriverForClient.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	defer netINodeDriverForClient.ReleaseNetINode(uNetINode)
	assert.NoError(t, err)

	writeData := make([]byte, 73)
	writeData[3] = 1
	writeData[7] = 2
	writeData[8] = 3
	writeData[33] = 4
	writeData[60] = 5
	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, 612))
	assert.NoError(t, netINodeDriverForClient.Sync(uNetINode))

	var maxWriteTimes int = 128
	for i = 0; i < maxWriteTimes; i++ {
		assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, uint64(netBlockCap*600+8*i)))
	}

	readData := make([]byte, 73)
	_, err = netINodeDriverForClient.PReadWithMem(uNetINode, readData, 612)
	assert.NoError(t, err)
	assert.Equal(t, writeData, readData)

	time.Sleep(time.Microsecond * 600)
	for i = 0; i < len(solodnSRPCListenAddrs); i++ {
		assert.NoError(t, solodns[i].Close())
	}
	assert.NoError(t, solonnIns.Close())
}
