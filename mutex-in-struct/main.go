package main

import (
	"github.com/go-auxiliaries/perf-geek/mutex-in-struct/code"
)

func main() {
	mut := (*code.Mutex)(nil).Init()
	mut.Read()
	nomut := (*code.NoMutex)(nil).Init()
	nomut.Read()

	go func() {
		mut.Read()
	}()

	go func() {
		nomut.Read()
	}()
}
