package code

import (
	"sync"
	"sync/atomic"
)

type MutexStruct struct {
	mu sync.Mutex
}

type EmbedMutex struct {
	val  atomic.Pointer[State]
	attr MutexStruct
}

func (d *EmbedMutex) Init() *EmbedMutex {
	tmp := EmbedMutex{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *EmbedMutex) Read() State {
	return *d.val.Load()
}

func (d EmbedMutex) ReadValueRcvr() State {
	return *d.val.Load()
}

func (d *EmbedMutex) Write(st State) {
	d.val.Store(&st)
}
