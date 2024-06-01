package code

import (
	"sync"
	"sync/atomic"
)

type PtrMutex struct {
	val atomic.Pointer[State]
	mut *sync.RWMutex
}

func (d *PtrMutex) Init() *PtrMutex {
	tmp := PtrMutex{}
	tmp.val.Store(&State{})
	tmp.mut = &sync.RWMutex{}
	return &tmp
}

func (d *PtrMutex) Read() State {
	return *d.val.Load()
}

func (d PtrMutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *PtrMutex) Write(st State) {
	d.val.Store(&st)
}
