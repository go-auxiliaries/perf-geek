package code

import (
	"sync"
	"sync/atomic"
)

type StateAtomicFixed struct {
	val           atomic.Pointer[State]
	_             [32]byte // To split atomic.Pointer and sync.RWMutex apart to avoid performance degradation https://github.com/golang/go/issues/67764
	exclusiveLock sync.RWMutex
}

func (d *StateAtomicFixed) Init() *StateAtomicFixed {
	return &StateAtomicFixed{}
}

func (d *StateAtomicFixed) runAtomically(body func(State) *State) {
	x := 0
	for {
		d.exclusiveLock.RLock()
		orig := d.val.Load()
		var newVal *State
		if orig == nil {
			newVal = body(State{})
		} else {
			newVal = body(*orig)
		}
		if newVal == nil {
			d.exclusiveLock.RUnlock()
			return
		}
		if d.val.CompareAndSwap(orig, newVal) {
			d.exclusiveLock.RUnlock()
			return
		}
		d.exclusiveLock.RUnlock()
		x++
		if x > 3 {
			d.exclusiveLock.Lock()
			if orig == nil {
				newVal = body(State{})
			} else {
				newVal = body(*orig)
			}
			if newVal == nil {
				d.exclusiveLock.Unlock()
				return
			}
			d.val.Store(newVal)
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
	val := d.val.Load()
	if val == nil {
		return State{}
	}
	return *val
}
