## Problem

Investigating performance problem for `state-consistency` I have found that **fact of existing** a `sync.Mutex` on the 
structure leads to 3 times read performance degradation for `atomic.Pointer.Load` in the same structure.

go version: `go version go1.22.3 linux/amd64`

Issue on golang repo: https://github.com/golang/go/issues/67764

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

### 5. Mutex atomic.Pointer and integer in the struct

```golang
type WithIntMutex struct {
    mut sync.RWMutex
    val atomic.Pointer[State]
    d   int
}

func (d *WithIntMutex) Read() State {
    return *d.val.Load()
}
```

### Results

```
goos: linux
goarch: amd64
pkg: github.com/go-auxiliaries/perf-geek/mutex-in-struct
cpu: 12th Gen Intel(R) Core(TM) i9-12900HK
BenchmarkTest
BenchmarkTest/single-thread
BenchmarkTest/single-thread/Mutex
BenchmarkTest/single-thread/Mutex/Read
BenchmarkTest/single-thread/Mutex/Read-12               505139860                2.332 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/Mutex/Write
BenchmarkTest/single-thread/Mutex/Write-12              27980103                41.06 ns/op           32 B/op          1 allocs/op
BenchmarkTest/single-thread/NoMutex
BenchmarkTest/single-thread/NoMutex/Read
BenchmarkTest/single-thread/NoMutex/Read-12             514682031                2.327 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/NoMutex/Write
BenchmarkTest/single-thread/NoMutex/Write-12            29626570                39.72 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread
BenchmarkTest/multi-thread/Mutex
BenchmarkTest/multi-thread/Mutex/Read
BenchmarkTest/multi-thread/Mutex/Read-12                584008728                1.775 ns/op           0 B/op          0 allocs/op
BenchmarkTest/multi-thread/Mutex/Write
BenchmarkTest/multi-thread/Mutex/Write-12               25739450                44.97 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/NoMutex
BenchmarkTest/multi-thread/NoMutex/Read
BenchmarkTest/multi-thread/NoMutex/Read-12              1000000000               0.4073 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/NoMutex/Write
BenchmarkTest/multi-thread/NoMutex/Write-12             27261466                43.73 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/PtrMutex
BenchmarkTest/multi-thread/PtrMutex/Read
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4113 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Write
BenchmarkTest/multi-thread/PtrMutex/Write-12            26372058                45.09 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/LikeMutex
BenchmarkTest/multi-thread/LikeMutex/Read
BenchmarkTest/multi-thread/LikeMutex/Read-12            1000000000               0.4084 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Write
BenchmarkTest/multi-thread/LikeMutex/Write-12           26154640                45.26 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/EmbedMutex
BenchmarkTest/multi-thread/EmbedMutex/Read
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4093 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Write
BenchmarkTest/multi-thread/EmbedMutex/Write-12          26043236                45.26 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/WithIntMutex
BenchmarkTest/multi-thread/WithIntMutex/Read
BenchmarkTest/multi-thread/WithIntMutex/Read-12         1000000000               0.4017 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/WithIntMutex/Write
BenchmarkTest/multi-thread/WithIntMutex/Write-12        30475310                38.22 ns/op           32 B/op          1 allocs/op
PASS
ok      github.com/go-auxiliaries/perf-geek/mutex-in-struct     16.129s
```

### Summary

`sync.Pointer.Store` operation is not affected:
```
BenchmarkTest/multi-thread/Mutex/Write-12               25739450                44.97 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/NoMutex/Write-12             27261466                43.73 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Write-12            26372058                45.09 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Write-12           26154640                45.26 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Write-12          26043236                45.26 ns/op           32 B/op          1 allocs/op
BenchmarkTest/multi-thread/WithIntMutex/Write-12        30475310                38.22 ns/op           32 B/op          1 allocs/op
```

Single threaded performance as well:
```
BenchmarkTest/single-thread/Mutex/Read-12               505139860                2.332 ns/op           0 B/op          0 allocs/op
BenchmarkTest/single-thread/NoMutex/Read-12             514682031                2.327 ns/op           0 B/op          0 allocs/op
```

But `sync.Pointer.Load` operation went from `0.4073` to `1.775` for the sheer fact of existence of a `sync.Mutex` on the structure:
```
BenchmarkTest/multi-thread/Mutex/Read-12                584008728                1.775 ns/op           0 B/op          0 allocs/op
BenchmarkTest/multi-thread/NoMutex/Read-12              1000000000               0.4073 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4113 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/LikeMutex/Read-12            1000000000               0.4084 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4093 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/WithIntMutex/Read-12         1000000000               0.4017 ns/op          0 B/op          0 allocs/op
```

It is `4.357`x performance degradation.

But having `sync.Mutex` as a pointer or encapsulated into another structure does not create performance drop:
```
BenchmarkTest/multi-thread/PtrMutex/Read-12             1000000000               0.4113 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/EmbedMutex/Read-12           1000000000               0.4093 ns/op          0 B/op          0 allocs/op
BenchmarkTest/multi-thread/WithIntMutex/Read-12         1000000000               0.4017 ns/op          0 B/op          0 allocs/op
```

Weirdest thing out of all is that adding an `int` to the structure makes performance drop go away:
```
BenchmarkTest/multi-thread/WithIntMutex/Read-12         1000000000               0.4017 ns/op          0 B/op          0 allocs/op
```


#### ASM Code

```go
	go func() {
		mut.Read()
	}()
```

Compiles to:
```asm
	go func() {
  0x45d760		493b6610		CMPQ SP, 0x10(R14)	
  0x45d764		7627			JBE 0x45d78d		
  0x45d766		55			PUSHQ BP		
  0x45d767		4889e5			MOVQ SP, BP		
  0x45d76a		4883ec10		SUBQ $0x10, SP		
  0x45d76e		488b4208		MOVQ 0x8(DX), AX	
	return *d.val.Load()
  0x45d772		8400			TESTB AL, 0(AX)		
		nomut.Read()
  0x45d774		90			NOPL			
	return *d.val.Load()
  0x45d775		488d1d04c00200		LEAQ sync/atomic..dict.Pointer[github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.State](SB), BX														
  0x45d77c		0f1f4000		NOPL 0(AX)																									
  0x45d780		e8dbfdffff		CALL sync/atomic.(*Pointer[go.shape.struct { github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.fruits []string; github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.numberOfApples int }]).Load(SB)	
  0x45d785		8400			TESTB AL, 0(AX)																									
	}()
  0x45d787		4883c410		ADDQ $0x10, SP		
  0x45d78b		5d			POPQ BP			
  0x45d78c		c3			RET
```

```go
	go func() {
        nomut.Read()
    }()
```

Compiles to:
```go
	go func() {
  0x45d7a0		493b6610		CMPQ SP, 0x10(R14)	
  0x45d7a4		7627			JBE 0x45d7cd		
  0x45d7a6		55			PUSHQ BP		
  0x45d7a7		4889e5			MOVQ SP, BP		
  0x45d7aa		4883ec10		SUBQ $0x10, SP		
  0x45d7ae		488b4208		MOVQ 0x8(DX), AX	
	return *d.val.Load()
  0x45d7b2		8400			TESTB AL, 0(AX)		
  0x45d7b4		4883c018		ADDQ $0x18, AX		
		mut.Read()
  0x45d7b8		90			NOPL			
	return *d.val.Load()
  0x45d7b9		488d1dc0bf0200		LEAQ sync/atomic..dict.Pointer[github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.State](SB), BX														
  0x45d7c0		e89bfdffff		CALL sync/atomic.(*Pointer[go.shape.struct { github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.fruits []string; github.com/go-auxiliaries/perf-geek/mutex-in-struct/code.numberOfApples int }]).Load(SB)	
  0x45d7c5		8400			TESTB AL, 0(AX)																									
	}()
  0x45d7c7		4883c410		ADDQ $0x10, SP		
  0x45d7cb		5d			POPQ BP			
  0x45d7cc		c3			RET
```

So, no difference on the bytecode