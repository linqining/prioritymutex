package prioritymutex

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type PriorityMutex struct {
	l      sync.Mutex
	pCount atomic.Int64
}

func (p *PriorityMutex) PLock() {
	p.pCount.Add(1)
	p.l.Lock()
}

func (p *PriorityMutex) PUnlock() {
	p.l.Unlock()
	p.pCount.Add(-1)
}

func (p *PriorityMutex) Lock() {
	if p.pCount.Load() > 0 {
		runtime.Gosched()
	}
	p.l.Lock()
}

func (p *PriorityMutex) Unlock() {
	p.l.Unlock()
}
