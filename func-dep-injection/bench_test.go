package func_dep_injection

import (
	"github.com/go-auxiliaries/perf-geek/func-dep-injection/code"
	"testing"
)

type FuncInterface struct{}

func (FuncInterface) CB(ch rune) int {
	if ch == 'l' {
		return 1
	}
	return 0
}

type FuncInterfaceFunc func(rune) int

func (fn FuncInterfaceFunc) CB(ch rune) int {
	return fn(ch)
}

func Benchmark(b *testing.B) {
	b.Run("FuncInterface", func(b *testing.B) {
		cb := FuncInterface{}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			code.CountFuncInterface("hello", cb)
		}
	})
	b.Run("FuncInterfaceFunc", func(b *testing.B) {
		cb := FuncInterfaceFunc(func(ch rune) int {
			if ch == 'l' {
				return 1
			}
			return 0
		})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			code.CountFuncInterface("hello", cb)
		}
	})
	b.Run("FuncInterfaceFunc", func(b *testing.B) {
		cb := FuncInterfaceFunc(func(ch rune) int {
			if ch == 'l' {
				return 1
			}
			return 0
		})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			code.CountFuncInterface("hello", cb)
		}
	})
}
