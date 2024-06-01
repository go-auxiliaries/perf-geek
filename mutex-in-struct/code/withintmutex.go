package code

import (
	"sync"
	"sync/atomic"
)

type WithIntMutex struct {
	mut sync.RWMutex
	val atomic.Pointer[State]
	d   int
}

func (d *WithIntMutex) Init() *WithIntMutex {
	tmp := WithIntMutex{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *WithIntMutex) Read() State {
	return *d.val.Load()
}

func (d WithIntMutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *WithIntMutex) Write(st State) {
	d.val.Store(&st)
}
