package netstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetBlockPool struct {
	options NetBlockPoolOptions
	driver  *NetBlockDriver
	pool    offheap.RawObjectPool
}

func (p *NetBlockPool) Init(options NetBlockPoolOptions,
	driver *NetBlockDriver) error {
	var err error

	p.options = options
	p.driver = driver

	err = p.driver.offheapDriver.InitRawObjectPool(&p.pool,
		int(types.NetBlockStructSize), p.options.RawChunksLimit,
		p.RawChunkPoolInvokePrepareNewRawChunk, p.RawChunkPoolInvokeReleaseRawChunk)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockPool) RawChunkPoolInvokeReleaseRawChunk() {
	panic("not support")
}

func (p *NetBlockPool) RawChunkPoolInvokePrepareNewRawChunk(uRawChunk uintptr) {
}

// MustGetNetBlock get or init a netBlock
func (p *NetBlockPool) MustGetBlock(netBlockID types.PtrBindIndex) (types.NetBlockUintptr, bool) {
	u, loaded := p.pool.MustGetRawObject(netBlockID)
	uNetBlock := (types.NetBlockUintptr)(u)
	return uNetBlock, loaded
}
