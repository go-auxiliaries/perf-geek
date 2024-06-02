package code

import (
	"sync"
	"sync/atomic"
)

type StateAtomicFixed struct {
	val           atomic.Pointer[State]
	exclusiveLock *sync.RWMutex // it is pointer because of https://github.com/golang/go/issues/67764
}

func (d *StateAtomicFixed) Init() *StateAtomicFixed {
	tmp := StateAtomicFixed{}
	tmp.val.Store(&State{})
	tmp.exclusiveLock = &sync.RWMutex{}
	return &tmp
}

func (d *StateAtomicFixed) runAtomically(body func(State) *State) {
	x := 0
	for {
		d.exclusiveLock.RLock()
		orig := d.val.Load()
		if d.val.CompareAndSwap(orig, body(*orig)) {
			d.exclusiveLock.RUnlock()
			return
		}
		d.exclusiveLock.RUnlock()
		x++
		if x > 3 {
			d.exclusiveLock.Lock()
			d.val.Store(body(*orig))
			d.exclusiveLock.Unlock()
			return
		}
	}
}

func (d *StateAtomicFixed) AddFruit(fruit string) {
	d.runAtomically(func(old State) *State {
		old.fruits = append(old.fruits, fruit)
		if fruit == "apple" {
			old.numberOfApples++
		}
		return &old
	})
}

func (d *StateAtomicFixed) GetState() State {
	return *d.val.Load()
}
