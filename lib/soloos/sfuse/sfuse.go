package sfuse

import (
	"soloos/sdfs/libsdfs"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type Server struct {
	options Options

	SFuseFs      SFuseFs
	MountOpts    fuse.MountOptions
	FuseServer   *fuse.Server
	FuseConnOpts *nodefs.Options
}

func (p *Server) Init(options Options, clientDriver *libsdfs.ClientDriver) error {
	var err error
	p.options = options

	err = p.SFuseFs.InitBySFuse(options, clientDriver)
	if err != nil {
		return err
	}

	p.MountOpts.AllowOther = true
	// p.MountOpts.MaxWrite = 0
	p.MountOpts.Name = FuseName
	p.MountOpts.Options = append(p.MountOpts.Options, "default_permissions")

	p.FuseServer, err = fuse.NewServer(&p.SFuseFs,
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
	err = p.SFuseFs.Close()
	if err != nil {
		return err
	}
	return nil
}
