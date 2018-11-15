package memstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type INodePool struct {
	options    INodePoolOptions
	inodeDriver *INodeDriver
	pool       offheap.RawObjectPool
}

func (p *INodePool) Init(options INodePoolOptions,
	inodeDriver *INodeDriver) error {
	var err error

	p.options = options
	p.inodeDriver = inodeDriver

	err = p.inodeDriver.offheapDriver.InitRawObjectPool(&p.pool,
		int(types.INodeStructSize), p.options.RawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	return nil
}

func (p *INodePool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *INodePool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetINode get or init a inodeblock
func (p *INodePool) MustGetINode(inodeID types.INodeID) (types.INodeUintptr, bool) {
	u, loaded := p.pool.MustGetRawObject(inodeID)
	uINode := (types.INodeUintptr)(u)
	return uINode, loaded
}
