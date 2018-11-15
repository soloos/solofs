package netstg

type NetBlockDriverOptions struct {
	NetBlockPoolOptions NetBlockPoolOptions
}

type NetBlockPoolOptions struct {
	RawChunksLimit int32
}
