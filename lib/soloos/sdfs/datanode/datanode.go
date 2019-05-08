package datanode

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
	"soloos/sdfs/localfs"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

type DataNode struct {
	*soloosbase.SoloOSEnv
	peerID  snettypes.PeerID
	metaStg *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	nameNodeClient api.NameNodeClient

	localFs        localfs.LocalFs
	uLocalDiskPeer snettypes.PeerUintptr

	srpcServerListenAddr string
	srpcServerServeAddr  string
	SRPCServer           DataNodeSRPCServer
}

func (p *DataNode) Init(soloOSEnv *soloosbase.SoloOSEnv,
	options DataNodeOptions,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var (
		uNameNodePeer snettypes.PeerUintptr
		err           error
	)

	p.SoloOSEnv = soloOSEnv
	p.peerID = options.PeerID

	p.metaStg = metaStg
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodeDriver = netINodeDriver

	uNameNodePeer, _ = p.SNetDriver.MustGetPeer(&options.NameNodePeerID, options.NameNodeSRPCServer,
		sdfsapitypes.DefaultSDFSRPCProtocol)
	err = p.nameNodeClient.Init(p.SoloOSEnv, uNameNodePeer)
	if err != nil {
		return err
	}

	p.srpcServerListenAddr = options.SrpcServerListenAddr
	p.srpcServerServeAddr = options.SrpcServerServeAddr
	err = p.SRPCServer.Init(p, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	err = p.localFs.Init(options.LocalFsRoot)
	if err != nil {
		return err
	}
	p.uLocalDiskPeer, _ = p.SNetDriver.MustGetPeer(&p.peerID, "", snettypes.ProtocolDisk)

	p.netBlockDriver.SetPReadMemBlockWithDisk(p.localFs.PReadMemBlockWithDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFs.UploadMemBlockWithDisk)
	p.netBlockDriver.SetHelper(&p.nameNodeClient, p.netBlockDriver.PrepareNetBlockMetaDataWithFanout)

	p.netINodeDriver.SetHelper(nil,
		p.metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.metaStg.PrepareNetINodeMetaDataWithStorDB,
		p.metaStg.NetINodeCommitSizeInDB,
	)

	return nil
}

func (p *DataNode) Serve() error {
	var err error
	err = p.nameNodeClient.RegisterDataNode(p.peerID, p.srpcServerServeAddr)
	if err != nil {
		return err
	}

	return p.SRPCServer.Serve()
}

func (p *DataNode) Close() error {
	return p.SRPCServer.Close()
}
