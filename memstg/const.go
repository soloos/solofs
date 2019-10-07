package memstg

const (
	MaxInt    = 1<<31 - 1
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int(MaxUint64 >> 1)
	MinInt64  = -MaxInt - 1

	TRUE  = 1
	FALSE = 0
)
