## Problem

Investigating performance problem for `state-consistency` I have found that **fact of existing** a `sync.Mutex` on the 
structure leads to 3 times read performance degradation for `atomic.Pointer.Load` in the same structure.

go version: `go version go1.22.3 linux/amd64`

## Solutions

### 1. Mutex on the structure (problematic case)

```golang
type Mutex struct {
    val           atomic.Pointer[State]
    mut sync.RWMutex
}

func (d *Mutex) Read() State {
    return *d.val.Load()
}
```

### 2. No Mutex on the structure

```golang
type NoMutex struct {
    val           atomic.Pointer[State]
}

func (d *NoMutex) Read() State {
    return *d.val.Load()
}
```

### 3. Mutex pointer on the structure

```golang
type PtrMutex struct {
    val           atomic.Pointer[State]
    mut *sync.RWMutex
}

func (d *PtrMutex) Read() State {
    return *d.val.Load()
}
```

### 4. Mutex embedded into the structure

```golang
type MutexStruct struct {
    mu sync.Mutex
}

type EmbedMutex struct {
    val  atomic.Pointer[State]
    attr MutexStruct
}

func (d *EmbedMutex) Read() State {
    return *d.val.Load()
}
```

### 5. Mutex-like structure

```golang
type MutexLike struct {
	state int32
	sema  uint32
}

type LikeMutex struct {
	val atomic.Pointer[State]
	mut MutexLike
}

func (d *LikeMutex) Read() State {
    return *d.val.Load()
}
```

### Results

```
goos: linux
goarch: amd64
pkg: github.com/go-auxiliaries/benchs/async/mutex-in-struct
cpu: 12th Gen Intel(R) Core(TM) i9-12900HK
BenchmarkTest
BenchmarkTest/single-thread
BenchmarkTest/single-thread/Mutex
BenchmarkTest/single-thread/Mutex/Read
BenchmarkTest/single-thread/Mutex/Read-12               521379552                2.215 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/Mutex/Write
BenchmarkTest/single-thread/Mutex/Write-12              29655799                40.09 ns/op           32 B/op          1 allocs/op
BenchmarkTest/single-thread/NoMutex
BenchmarkTest/single-thread/NoMutex/Read
BenchmarkTest/single-thread/NoMutex/Read-12             534564370                2.205 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/NoMutex/Write
BenchmarkTest/single-thread/NoMutex/Write-12            29692969                40.08 ns/op           32 B/op          1 allocs/op
BenchmarkTest/single-thread/PtrMutex
BenchmarkTest/single-thread/PtrMutex/Read
BenchmarkTest/single-thread/PtrMutex/Read-12            494957496                2.239 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/PtrMutex/Write
BenchmarkTest/single-thread/PtrMutex/Write-12           28901656                39.91 ns/op           32 B/op          1 allocs/op
BenchmarkTest/single-thread/LikeMutex
BenchmarkTest/single-thread/LikeMutex/Read
BenchmarkTest/single-thread/LikeMutex/Read-12           533799769                2.206 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/LikeMutex/Write
BenchmarkTest/single-thread/LikeMutex/Write-12          29357671                39.97 ns/op           32 B/op          1 allocs/op
BenchmarkTest/single-thread/EmbedMutex
BenchmarkTest/single-thread/EmbedMutex/Read
BenchmarkTest/single-thread/EmbedMutex/Read-12          499452915                2.302 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/EmbedMutex/Write
BenchmarkTest/single-thread/EmbedMutex/Write-12         29465785                40.98 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread
BenchmarkTest/multi-thread/Mutex
BenchmarkTest/multi-thread/Mutex/Read
BenchmarkTest/multi-thread/Mutex/Read-12                873189706                2.348 ns/op           0 B/op          0 allocs/op
BenchmarkTest/multi-thread/Mutex/Write
BenchmarkTest/multi-thread/Mutex/Write-12               24801372                45.66 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/NoMutex
BenchmarkTest/multi-thread/NoMutex/Read
BenchmarkTest/multi-thread/NoMutex/Read-12              1000000000               0.3980 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/NoMutex/Write
BenchmarkTest/multi-thread/NoMutex/Write-12             26100777                45.45 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/PtrMutex
BenchmarkTest/multi-thread/PtrMutex/Read
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4107 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Write
BenchmarkTest/multi-thread/PtrMutex/Write-12            26942925                43.71 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/LikeMutex
BenchmarkTest/multi-thread/LikeMutex/Read
BenchmarkTest/multi-thread/LikeMutex/Read-12            1000000000               0.4015 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Write
BenchmarkTest/multi-thread/LikeMutex/Write-12           25786527                45.61 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/EmbedMutex
BenchmarkTest/multi-thread/EmbedMutex/Read
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4036 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Write
BenchmarkTest/multi-thread/EmbedMutex/Write-12          25457143                45.74 ns/op           32 B/op          1 allocs/op
PASS
ok      github.com/go-auxiliaries/benchs/async/mutex-in-struct  23.979s
```

### Summary

`sync.Pointer.Store` operation is not affected:
```
BenchmarkTest/multi-thread/Mutex/Write-12               24801372                45.66 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/NoMutex/Write-12             26100777                45.45 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Write-12            26942925                43.71 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Write-12           25786527                45.61 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Write-12          25457143                45.74 ns/op           32 B/op          1 allocs/op
```

Single threaded performance as well:
```
BenchmarkTest/single-thread/Mutex/Read-12               521379552                2.215 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/NoMutex/Read-12             534564370                2.205 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/PtrMutex/Read-12            494957496                2.239 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/LikeMutex/Read-12           533799769                2.206 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/EmbedMutex/Read-12          499452915                2.302 ns/op           0 B/op          0 allocs/op
```

But `sync.Pointer.Load` operation went from `0.4107` to `2.348` for the sheer fact of existence of a `sync.Mutex` on the structure:
```
BenchmarkTest/multi-thread/Mutex/Read-12                873189706                2.348 ns/op           0 B/op          0 allocs/op
BenchmarkTest/multi-thread/NoMutex/Read-12              1000000000               0.3980 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4107 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Read-12            1000000000               0.4015 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4036 ns/op          0 B/op          0 allocs/op
```

It is `5.717` performance degradation.

But having `sync.Mutex` as a pointer or encapsulated into another structure does not create performance drop:
```
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4107 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4036 ns/op          0 B/op          0 allocs/op
```
