package code

import (
	"sync/atomic"
)

type NoMutex struct {
	val atomic.Pointer[State]
}

func (d *NoMutex) Init() *NoMutex {
	tmp := NoMutex{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *NoMutex) Read() State {
	return *d.val.Load()
}

func (d NoMutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *NoMutex) Write(st State) {
	d.val.Store(&st)
}
