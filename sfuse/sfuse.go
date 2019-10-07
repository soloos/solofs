package sfuse

import (
	"os"
	"soloos/common/go-fuse/fuse"
	"soloos/common/solofstypes"
	"soloos/solofs/solofssdk"
)

type Server struct {
	options Options

	Client     solofssdk.Client
	MountOpts  fuse.MountOptions
	FuseServer *fuse.Server
}

func (p *Server) Init(options Options,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	defaultMemBlocksLimit int32,
	clientDriver *solofssdk.ClientDriver) error {
	var err error
	p.options = options

	os.MkdirAll(options.MountPoint, 0777)

	err = clientDriver.InitClient(&p.Client,
		solofstypes.NameSpaceID(options.NameSpaceID),
		defaultNetBlockCap, defaultMemBlockCap, defaultMemBlocksLimit)
	if err != nil {
		return err
	}

	p.MountOpts.AllowOther = true
	// p.MountOpts.MaxWrite = 0
	p.MountOpts.Name = p.Client.GetPosixFs().String() + "-fuse"
	p.MountOpts.Options = append(p.MountOpts.Options, "default_permissions")

	p.FuseServer, err = fuse.NewServer(p.Client.GetPosixFs(),
		p.options.MountPoint,
		&p.MountOpts)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) Serve() error {
	p.FuseServer.Serve()
	return nil
}

func (p *Server) Close() error {
	var err error
	err = p.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
