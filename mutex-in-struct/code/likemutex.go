package code

import (
	"sync/atomic"
)

type MutexLike struct {
	state int32
	sema  uint32
}

type LikeMutex struct {
	val atomic.Pointer[State]
	mut MutexLike
}

func (d *LikeMutex) Init() *LikeMutex {
	tmp := LikeMutex{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *LikeMutex) Read() State {
	return *d.val.Load()
}

func (d LikeMutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *LikeMutex) Write(st State) {
	d.val.Store(&st)
}
