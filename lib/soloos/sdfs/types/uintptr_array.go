package types

type UintptrArray8 struct {
	Arr [8]uintptr
	Len int
}

func (p *UintptrArray8) Append(value uintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type UintptrArray16 struct {
	Arr [16]uintptr
	Len int
}

func (p *UintptrArray16) Append(value uintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type UintptrArray32 struct {
	Arr [32]uintptr
	Len int
}

func (p *UintptrArray32) Append(value uintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}

type UintptrArray64 struct {
	Arr [64]uintptr
	Len int
}

func (p *UintptrArray64) Append(value uintptr) {
	p.Arr[p.Len] = value
	p.Len += 1
}
