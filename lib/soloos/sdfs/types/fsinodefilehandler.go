package types

type FsINodeFileHandler struct {
	FsINodeID      FsINodeID
	AppendPosition uint64
	ReadPosition   uint64
}

func (p *FsINodeFileHandler) Reset() {
	p.AppendPosition = 0
	p.ReadPosition = 0
}
