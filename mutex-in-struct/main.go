package main

import (
	"github.com/go-auxiliaries/perf-geek/mutex-in-struct/code"
)

func main() {
	tmp := code.Mutex{}
	tmp.Read()
	tmp2 := code.NoMutex{}
	tmp2.Read()

	go func() {
		tmp := code.Mutex{}
		tmp.Read()
	}()

	go func() {
		tmp := code.NoMutex{}
		tmp.Read()
	}()

}
