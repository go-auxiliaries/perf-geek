package main_test

import (
	"github.com:go-auxiliaries/perf-geek/mutex-in-struct/code"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	b.Run("single-thread", func(b *testing.B) {
		b.Run("Mutex", func(b *testing.B) {
			runSingleThreadTest[*code.Mutex](b)
		})
		b.Run("NoMutex", func(b *testing.B) {
			runSingleThreadTest[*code.NoMutex](b)
		})
		b.Run("PtrMutex", func(b *testing.B) {
			runSingleThreadTest[*code.PtrMutex](b)
		})
		b.Run("LikeMutex", func(b *testing.B) {
			runSingleThreadTest[*code.LikeMutex](b)
		})
		b.Run("EmbedMutex", func(b *testing.B) {
			runSingleThreadTest[*code.EmbedMutex](b)
		})
	})

	b.Run("multi-thread", func(b *testing.B) {
		b.Run("Mutex", func(b *testing.B) {
			runMultiThreadTest[*code.Mutex](b)
		})
		b.Run("NoMutex", func(b *testing.B) {
			runMultiThreadTest[*code.NoMutex](b)
		})
		b.Run("PtrMutex", func(b *testing.B) {
			runMultiThreadTest[*code.PtrMutex](b)
		})
		b.Run("LikeMutex", func(b *testing.B) {
			runMultiThreadTest[*code.LikeMutex](b)
		})
		b.Run("EmbedMutex", func(b *testing.B) {
			runMultiThreadTest[*code.EmbedMutex](b)
		})
	})
}

type objectUnderTest[T any] interface {
	Init() T
	Read() code.State
	ReadValueRcvr() code.State
	Write(state code.State)
}

func runMultiThreadTest[O objectUnderTest[O]](b *testing.B) {
	o := O.Init(*new(O))
	b.Run("Read", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = o.Read()
			}
		})
	})
	o = O.Init(*new(O))
	b.Run("Write", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				o.Write(code.State{})
			}
		})
	})
}

func runSingleThreadTest[O objectUnderTest[O]](b *testing.B) {
	o := O.Init(*new(O))
	b.Run("Read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = o.Read()
		}
	})
	o = O.Init(*new(O))
	b.Run("Write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			o.Write(code.State{})
		}
	})
}
