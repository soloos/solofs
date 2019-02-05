package sfuse

import (
	"os"
	"soloos/sdfs/libsdfs"
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type Server struct {
	options Options

	Client       libsdfs.Client
	MountOpts    fuse.MountOptions
	FuseServer   *fuse.Server
	FuseConnOpts *nodefs.Options
}

func (p *Server) Init(options Options,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	clientDriver *libsdfs.ClientDriver) error {
	var err error
	p.options = options

	os.MkdirAll(options.MountPoint, 0777)

	err = clientDriver.InitClient(&p.Client, defaultNetBlockCap, defaultMemBlockCap)
	if err != nil {
		return err
	}

	p.MountOpts.AllowOther = true
	// p.MountOpts.MaxWrite = 0
	p.MountOpts.Name = types.FuseName
	p.MountOpts.Options = append(p.MountOpts.Options, "default_permissions")

	p.FuseServer, err = fuse.NewServer(&p.Client.MemDirTreeStg,
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
