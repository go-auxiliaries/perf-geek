package code

import (
	"sync"
	"sync/atomic"
)

type Mutex struct {
	val atomic.Pointer[State]
	mut sync.RWMutex
}

func (d *Mutex) Init() *Mutex {
	tmp := Mutex{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *Mutex) Read() State {
	return *d.val.Load()
}

func (d Mutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *Mutex) Write(st State) {
	d.val.Store(&st)
}
