package test_from_issue

import (
	"sync"
	"sync/atomic"
	"testing"
)

type State struct {
	l []string
	i int
}

type State2 struct {
	l []string
	i int
	k int
}

type Mutex[T any] struct {
	mut sync.RWMutex
	val atomic.Pointer[T]
}

func (d *Mutex[T]) Read() T {
	return *d.val.Load()
}

type MutexWithInt[T any] struct {
	mut sync.RWMutex
	val atomic.Pointer[T]
	i   int
}

func (d *MutexWithInt[T]) Read() T {
	return *d.val.Load()
}

func BenchmarkMutexWith_State(b *testing.B) {
	o := Mutex[State]{}
	o.val.Store(&State{})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = o.Read()
		}
	})
}

func BenchmarkMutexWith_State2(b *testing.B) {
	o := Mutex[State2]{}
	o.val.Store(&State2{})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = o.Read()
		}
	})
}

func BenchmarkMutexWithInt_State(b *testing.B) {
	o := MutexWithInt[State]{}
	o.val.Store(&State{})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = o.Read()
		}
	})
}

func BenchmarkMutexWithInt_State2(b *testing.B) {
	o := MutexWithInt[State2]{}
	o.val.Store(&State2{})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = o.Read()
		}
	})
}
