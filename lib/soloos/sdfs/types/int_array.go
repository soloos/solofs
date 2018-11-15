package types

type Int64Array8 struct {
	Arr [8]int64
	Len int
}

func (p *Int64Array8) Append(value int64) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type Int64Array16 struct {
	Arr [16]int64
	Len int
}

func (p *Int64Array16) Append(value int64) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type Int64Array32 struct {
	Arr [32]int64
	Len int
}

func (p *Int64Array32) Append(value int64) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type Int64Array64 struct {
	Arr [64]int64
	Len int
}

func (p *Int64Array64) Append(value int64) {
	p.Arr[p.Len] = value
	p.Len += 1
}
