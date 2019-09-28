package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFs) SetNetINodeBlockPlacement(netINodeID solofsapitypes.NetINodeID,
	policy solofsapitypes.MemBlockPlacementPolicy) error {
	var (
		uNetINode solofsapitypes.NetINodeUintptr
		err       error
	)
	uNetINode, err = p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
	defer p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	if uNetINode == 0 {
		return solofsapitypes.ErrNetINodeNotExists
	}

	uNetINode.Ptr().MemBlockPlacementPolicy = policy

	return nil
}
