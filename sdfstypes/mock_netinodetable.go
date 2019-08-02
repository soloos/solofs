package sdfstypes

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
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
		int(sdfsapitypes.NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockNetINodeTable) MustGetNetINode(netINodeID sdfsapitypes.NetINodeID) (sdfsapitypes.NetINodeUintptr, bool) {
	uObject, afterSetNewObj := p.table.MustGetObject(netINodeID)
	var loaded = afterSetNewObj == nil
	if afterSetNewObj != nil {
		afterSetNewObj()
	}
	uNetINode := (sdfsapitypes.NetINodeUintptr)(uObject)
	return uNetINode, loaded
}

func (p *MockNetINodeTable) AllocNetINode(netBlockCap, memBlockCap int) sdfsapitypes.NetINodeUintptr {
	var netINodeID sdfsapitypes.NetINodeID
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
