package types

import "sync"

// TODO use lru
type HotPool struct {
	mutex   sync.Mutex
	tmppool []interface{}
}

func (p *HotPool) Init() error {
	return nil
}

func (p *HotPool) Pop() interface{} {
	p.mutex.Lock()
	last := len(p.tmppool) - 1
	x := p.tmppool[last]
	p.tmppool = p.tmppool[:last]
	p.mutex.Unlock()
	return x
}

func (p *HotPool) Put(x interface{}) {
	p.mutex.Lock()
	p.tmppool = append(p.tmppool, x)
	p.mutex.Unlock()
}

func (p *HotPool) IteratorAndPop(itFunc func(x interface{}) (bool, interface{})) interface{} {
	var (
		isBreak  bool
		ret      interface{}
		popIndex int
	)
	p.mutex.Lock()
	for k, _ := range p.tmppool {
		isBreak, ret = itFunc(p.tmppool[k])
		if isBreak {
			popIndex = k
			break
		}
	}
	p.tmppool = append(p.tmppool[:popIndex], p.tmppool[popIndex+1:]...)
	p.mutex.Unlock()
	return ret
}
