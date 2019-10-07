package memstg

import (
	"soloos/common/solofstypes"
)

func (p *PosixFs) SetNetINodeBlockPlacement(netINodeID solofstypes.NetINodeID,
	policy solofstypes.MemBlockPlacementPolicy) error {
	var (
		uNetINode solofstypes.NetINodeUintptr
		err       error
	)
	uNetINode, err = p.MemStg.NetINodeDriver.GetNetINode(netINodeID)
	defer p.MemStg.NetINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	if uNetINode == 0 {
		return solofstypes.ErrNetINodeNotExists
	}

	uNetINode.Ptr().MemBlockPlacementPolicy = policy

	return nil
}
