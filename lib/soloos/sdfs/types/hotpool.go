package types

import (
	"sync"
)

// TODO use lru
type HotPool struct {
	mutex   sync.Mutex
	tmppool []uintptr
}

func (p *HotPool) Init() error {
	return nil
}

func (p *HotPool) Pop() uintptr {
	p.mutex.Lock()
	last := len(p.tmppool) - 1
	x := p.tmppool[last]
	p.tmppool = p.tmppool[:last]
	p.mutex.Unlock()
	return x
}

func (p *HotPool) Put(x uintptr) {
	p.mutex.Lock()
	p.tmppool = append(p.tmppool, x)
	p.mutex.Unlock()
}

func (p *HotPool) IteratorAndPop(itFunc func(x uintptr) (bool, uintptr)) uintptr {
	var (
		isBreak  bool
		ret      uintptr
		popIndex int = -1
	)
	p.mutex.Lock()
	for k, _ := range p.tmppool {
		isBreak, ret = itFunc(p.tmppool[k])
		if isBreak {
			popIndex = k
			break
		}
	}

	if popIndex >= 0 {
		if len(p.tmppool) > 1 {
			p.tmppool = append(p.tmppool[:popIndex], p.tmppool[popIndex+1:]...)
		} else {
			p.tmppool = p.tmppool[:0]
		}
	}
	p.mutex.Unlock()
	return ret
}
