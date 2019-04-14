package memstg

type MemBlockDriverOptions struct {
	MemBlockTableOptionsList []MemBlockTableOptions
}

type MemBlockTableOptions struct {
	ObjectSize   int
	ObjectsLimit int32
}
