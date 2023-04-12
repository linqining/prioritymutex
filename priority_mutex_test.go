package prioritymutex

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func priority(ctx context.Context, p *PriorityMutex, lockFlag chan struct{}) {
	p.PLock()
	defer p.Unlock()
	lockFlag <- struct{}{}
	<-ctx.Done()
}

var flag int32

func normal(p *PriorityMutex, lockFlag chan struct{}) {
	p.Lock()
	defer p.Unlock()
	atomic.AddInt32(&flag, 1)
	lockFlag <- struct{}{}
}

func TestPriorityMutex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	p := &PriorityMutex{}
	pLockFlag := make(chan struct{}, 1)
	go priority(ctx, p, pLockFlag)
	<-pLockFlag

	lockFlag := make(chan struct{}, 1)
	go normal(p, lockFlag)

	timeOutCtx, timeCancel := context.WithTimeout(context.TODO(), time.Second)
	defer timeCancel()
	select {
	case <-lockFlag:
		t.Fatal("unreachable")
	case <-timeOutCtx.Done():
	}
	if atomic.LoadInt32(&flag) != 0 {
		t.Fatal("flag change,it should be block and not change")
	}
	cancel()
	time.Sleep(time.Second)
	<-lockFlag
	if atomic.LoadInt32(&flag) != 1 {
		t.Fatal("flag not change")
	}
}
