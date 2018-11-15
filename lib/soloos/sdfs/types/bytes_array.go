package types

type BytesPtrArray8 struct {
	Arr [8]BytesUintptr
	Len int
}

func (p *BytesPtrArray8) Append(value BytesUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type BytesPtrArray16 struct {
	Arr [16]BytesUintptr
	Len int
}

func (p *BytesPtrArray16) Append(value BytesUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type BytesPtrArray32 struct {
	Arr [32]BytesUintptr
	Len int
}

func (p *BytesPtrArray32) Append(value BytesUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type BytesPtrArray64 struct {
	Arr [64]BytesUintptr
	Len int
}

func (p *BytesPtrArray64) Append(value BytesUintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}
