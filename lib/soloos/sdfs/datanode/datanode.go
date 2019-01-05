package datanode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/localfs"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type DataNode struct {
	offheapDriver    *offheap.OffheapDriver
	snetDriver       *snet.NetDriver
	snetClientDriver *snet.ClientDriver
	metaStg          *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	nameNodeClient api.NameNodeClient

	peerID         snettypes.PeerID
	localFs        localfs.LocalFs
	uLocalDiskPeer snettypes.PeerUintptr

	srpcServerListenAddr string
	srpcServerServeAddr  string
	SRPCServer           DataNodeSRPCServer
}

func (p *DataNode) Init(offheapDriver *offheap.OffheapDriver,
	options DataNodeOptions,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var (
		uNameNodePeer snettypes.PeerUintptr
		err           error
	)

	p.offheapDriver = offheapDriver
	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.metaStg = metaStg
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodeDriver = netINodeDriver

	uNameNodePeer, _ = p.snetDriver.MustGetPeer(&options.NameNodePeerID, options.NameNodeSRPCServer,
		types.DefaultSDFSRPCProtocol)
	err = p.nameNodeClient.Init(p.snetClientDriver, uNameNodePeer)
	if err != nil {
		return err
	}

	p.srpcServerListenAddr = options.SrpcServerListenAddr
	p.srpcServerServeAddr = options.SrpcServerServeAddr
	err = p.SRPCServer.Init(p, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.peerID = options.PeerID
	err = p.localFs.Init(options.LocalFsRoot)
	if err != nil {
		return err
	}
	p.uLocalDiskPeer, _ = p.snetDriver.MustGetPeer(&p.peerID, "", snettypes.ProtocolDisk)

	p.netBlockDriver.SetPReadMemBlockWithDisk(p.localFs.PReadMemBlockWithDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFs.UploadMemBlockWithDisk)
	p.netBlockDriver.SetHelper(&p.nameNodeClient, p.netBlockDriver.PrepareNetBlockMetaDataWithFanout)

	p.netINodeDriver.SetHelper(nil,
		p.metaStg.PrepareNetINodeMetaDataOnlyLoadDB, p.metaStg.PrepareNetINodeMetaDataWithStorDB)

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
