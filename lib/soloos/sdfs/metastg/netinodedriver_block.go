package metastg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *NetINodeDriver) ChooseDataNodesForNewNetBlock(uNetINode types.NetINodeUintptr,
	backends *snettypes.PeerUintptrArray8) error {
	backends.Reset()
	for i := 0; i < 3; i++ {
		backends.Append(p.helper.ChooseOneDataNode())
	}
	return nil
}
