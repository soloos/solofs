package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *PosixFS) SetNetINodeBlockPlacement(netINodeID sdfsapitypes.NetINodeID,
	policy sdfsapitypes.MemBlockPlacementPolicy) error {
	var (
		uNetINode sdfsapitypes.NetINodeUintptr
		err       error
	)
	uNetINode, err = p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
	defer p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	if uNetINode == 0 {
		return sdfsapitypes.ErrNetINodeNotExists
	}

	uNetINode.Ptr().MemBlockPlacementPolicy = policy

	return nil
}
