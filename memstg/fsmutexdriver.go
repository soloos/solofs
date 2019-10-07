package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type FsMutexDriver struct {
	*soloosbase.SoloosEnv
	posixFs           *PosixFs
	INodeRWMutexTable offheap.LKVTableWithUint64
}

func (p *FsMutexDriver) Init(
	soloosEnv *soloosbase.SoloosEnv,
	posixFs *PosixFs,
) error {
	var err error
	p.SoloosEnv = soloosEnv
	p.posixFs = posixFs

	err = p.SoloosEnv.OffheapDriver.InitLKVTableWithUint64(&p.INodeRWMutexTable, "INodeRWMutex",
		int(solofsapitypes.INodeRWMutexStructSize), -1, offheap.DefaultKVTableSharedCount,
		nil)
	if err != nil {
		return err
	}

	return nil
}
