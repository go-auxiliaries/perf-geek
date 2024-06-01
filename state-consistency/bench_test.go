package state_consistency

import (
	"fmt"
	"github.com/go-auxiliaries/benchs/async/state-consistency/code"
	"runtime"
	"testing"
)

func BenchmarkDataSync(b *testing.B) {
	b.Run("StateMutex", func(b *testing.B) {
		runSingleThreadTest[*code.StateMutex](b)
		runMultiThreadTest[*code.StateMutex](b)
	})
	b.Run("StateAtomicFixed", func(b *testing.B) {
		runSingleThreadTest[*code.StateAtomicFixed](b)
		runMultiThreadTest[*code.StateAtomicFixed](b)
	})
	b.Run("StateAtomic", func(b *testing.B) {
		runSingleThreadTest[*code.StateAtomic](b)
		runMultiThreadTest[*code.StateAtomic](b)
	})
}

type objectUnderTest[T any] interface {
	Init() T
	GetState() code.State
	AddFruit(string)
}

func runSingleThreadTest[O objectUnderTest[O]](b *testing.B) {
	b.Run("single-thread", func(b *testing.B) {
		b.Run("WriteState", func(b *testing.B) {
			runtime.GC()
			o := O.Init(*new(O))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				o.AddFruit("somefruit")
			}
		})
		b.Run("ReadState", func(b *testing.B) {
			var val code.State
			runtime.GC()
			o := O.Init(*new(O))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				val = o.GetState()
			}
			_ = val
		})
	})
}

func runMultiThreadTest[O objectUnderTest[O]](b *testing.B) {
	b.Run("multi-thread", func(b *testing.B) {
		b.Run("WriteState", func(b *testing.B) {
			o := O.Init(*new(O))
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					o.AddFruit("somefruit")
				}
			})
		})
		b.Run("ReadState", func(b *testing.B) {
			o := O.Init(*new(O))
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = o.GetState()
				}
			})
		})

		for _, rates := range [][2]int{{1, 20}, {1, 9}, {3, 7}, {5, 5}, {7, 3}, {9, 1}} {
			b.Run(fmt.Sprintf("Mixed-%d/%d", rates[0], rates[1]), func(b *testing.B) {
				o := O.Init(*new(O))
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for range rates[0] {
							o.AddFruit("somefruit")
						}
						for range rates[1] {
							_ = o.GetState()
						}
					}
				})
			})
		}
	})
}
