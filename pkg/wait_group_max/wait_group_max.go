package wait_group_max

import (
	"sync"
	"sync/atomic"
	"time"
)

type WaitGroupMax struct {
	wg      sync.WaitGroup
	max     int
	current int32
}

func CreateWorkGroupMax(max int) *WaitGroupMax {
	return &WaitGroupMax{max: max, current: 0}
}

func (w *WaitGroupMax) Add(delta int) {
	for atomic.LoadInt32(&w.current) >= int32(w.max) {
		//TODO: Replace sleep
		time.Sleep(1)
	}
	atomic.StoreInt32(&w.current, atomic.LoadInt32(&w.current)+1)
	w.wg.Add(delta)
}

func (w *WaitGroupMax) Done() {
	atomic.StoreInt32(&w.current, atomic.LoadInt32(&w.current)-1)
	w.wg.Done()
}

func (w *WaitGroupMax) Wait() {
	w.wg.Wait()
}
