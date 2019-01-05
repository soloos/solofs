package libsdfs

import "soloos/sdfs/types"

func (p *Client) FileTableAlloc(fsINodeID types.FsINodeID) uint64 {
	var fdID = int(p.FileTableIndexs.Get())
	p.FileTableRWMutex.RLock()
	p.FileTable[fdID].FsINodeID = fsINodeID
	p.FileTableRWMutex.RUnlock()
	return uint64(fdID)

}

func (p *Client) FileTableAddAppendPosition(fdID uint64, delta int64) {
	p.FileTableRWMutex.RLock()
	p.FileTable[int(fdID)].AppendPosition += delta
	p.FileTableRWMutex.RUnlock()
	return
}

func (p *Client) FileTableAddReadPosition(fdID uint64, delta int64) {
	p.FileTableRWMutex.RLock()
	p.FileTable[int(fdID)].ReadPosition += delta
	p.FileTableRWMutex.RUnlock()
	return
}

func (p *Client) FileTableGet(fdID uint64) (ret types.FsINodeFileHandler) {
	p.FileTableRWMutex.RLock()
	ret = p.FileTable[int(fdID)]
	p.FileTableRWMutex.RUnlock()
	return
}

func (p *Client) FileTableRelease(fdID uint64) {
	p.FileTable[int(fdID)].Reset()
	p.FileTableIndexs.Put(uintptr(fdID))
}
