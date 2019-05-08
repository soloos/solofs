package types

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdbone/offheap"
)

type MockNetINodeTable struct {
	*soloosbase.SoloOSEnv
	table offheap.LKVTableWithBytes64
}

func (p *MockNetINodeTable) Init(soloOSEnv *soloosbase.SoloOSEnv) error {
	var err error
	p.SoloOSEnv = soloOSEnv

	err = p.OffheapDriver.InitLKVTableWithBytes64(&p.table, "MockNetINode",
		int(NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockNetINodeTable) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	uObject, afterSetNewObj := p.table.MustGetObject(netINodeID)
	var loaded = afterSetNewObj == nil
	if afterSetNewObj != nil {
		afterSetNewObj()
	}
	uNetINode := (NetINodeUintptr)(uObject)
	return uNetINode, loaded
}

func (p *MockNetINodeTable) AllocNetINode(netBlockCap, memBlockCap int) NetINodeUintptr {
	var netINodeID NetINodeID
	sdfsapitypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, _ := p.MustGetNetINode(netINodeID)
	uNetINode.Ptr().ID = netINodeID
	uNetINode.Ptr().NetBlockCap = netBlockCap
	uNetINode.Ptr().MemBlockCap = memBlockCap
	return uNetINode
}

func (p *MockNetINodeTable) ReleaseNetINode(uNetINode sdfsapitypes.NetINodeUintptr) {
	p.table.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}
