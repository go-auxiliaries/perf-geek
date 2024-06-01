package code

import (
	"sync/atomic"
)

type StateAtomic struct {
	val atomic.Pointer[State]
}

func (d *StateAtomic) Init() *StateAtomic {
	tmp := StateAtomic{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *StateAtomic) runAtomically(body func(State) *State) {
	res := false
	for !res {
		orig := d.val.Load()
		res = d.val.CompareAndSwap(orig, body(*orig))
	}
}

func (d *StateAtomic) AddFruit(fruit string) {
	d.runAtomically(func(old State) *State {
		old.fruits = append(old.fruits, fruit)
		if fruit == "apple" {
			old.numberOfApples++
		}
		return &old
	})
}

func (d *StateAtomic) GetState() State {
	return *d.val.Load()
}
