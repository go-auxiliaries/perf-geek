## Problem

You have some complex state you want to manipulate it, but you want every manipulation, even complex to be atomic.
If two operations are running at the same time, they won't interfere with each other.

## Solution

### 1. RW Lock (StateMutex)
```golang
type StateMutex struct {       
	state State
	mu    sync.RWMutex
}

func (d *StateMutex) Init() *StateMutex {
	return &StateMutex{}
}

func (d *StateMutex) AddFruit(fruit string) {
	d.mu.Lock()
	d.state.fruits = append(d.state.fruits, fruit)
	if fruit == "apple" {
		d.state.numberOfApples++
	}
	d.mu.Unlock()
}

func (d *StateMutex) GetState() State {
	d.mu.RLock()
	val := d.state
	d.mu.RUnlock()
	return val
}
```

### 2. Use atomic pointer to a state

```golang
type State struct {
	data atomic.Pointer[StateData]
}

type StateData struct {
    CustomerNames []string
    NumberOfJones int
    // ...
}

func (s *State) runAtomically(body func(StateData) *StateData) {
    res := true
    for res {
		orig := *s.data.Load()
		new := body(orig)
        res = !d.val.CompareAndSwap(orig, new))
    }
}

func (s *State) AddCustomer(name string) {
	s.runAtomically(func(d StateData) *StateData {
        d.CustomerNames = append(d.CustomerNames, name)
        d.NumberOfJones += strings.Count(name, "Jones")
        return &d
    })
}
```

### 3. Use atomic pointer to a state, with mutex to negate state submitting collisions

```golang
type StateAtomicFixed struct {
	val           atomic.Pointer[State]
	exclusiveLock sync.RWMutex
}

func (d *StateAtomicFixed) Init() *StateAtomicFixed {
	tmp := StateAtomicFixed{}
	tmp.val.Store(&State{})
	return &tmp
}

func (d *StateAtomicFixed) runAtomically(body func(State) *State) {
	x := 0
	for {
		d.exclusiveLock.RLock()
		orig := d.val.Load()
		if d.val.CompareAndSwap(orig, body(*orig)) {
			d.exclusiveLock.RUnlock()
			return
		}
		d.exclusiveLock.RUnlock()
		x++
		if x > 3 {
			d.exclusiveLock.Lock()
			d.val.Store(body(*orig))
			d.exclusiveLock.Unlock()
			return
		}
	}
}
```

## Results
```
goos: linux
goarch: amd64
pkg: github.com/go-auxiliaries/benchs/async/state-consistency
cpu: 12th Gen Intel(R) Core(TM) i9-12900HK
BenchmarkDataSync
BenchmarkDataSync/StateMutex
BenchmarkDataSync/StateMutex/single-thread
BenchmarkDataSync/StateMutex/single-thread/WriteState
BenchmarkDataSync/StateMutex/single-thread/WriteState-12             21621788                50.89 ns/op           93 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/single-thread/ReadState
BenchmarkDataSync/StateMutex/single-thread/ReadState-12              98048242                12.19 ns/op            0 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread
BenchmarkDataSync/StateMutex/multi-thread/WriteState
BenchmarkDataSync/StateMutex/multi-thread/WriteState-12              13214971               102.9 ns/op            97 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/ReadState
BenchmarkDataSync/StateMutex/multi-thread/ReadState-12               24539293                48.51 ns/op            0 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/Mixed-1/9
BenchmarkDataSync/StateMutex/multi-thread/Mixed-1/9-12                6304790               188.0 ns/op            83 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/Mixed-3/7
BenchmarkDataSync/StateMutex/multi-thread/Mixed-3/7-12                3736591               320.4 ns/op           276 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/Mixed-5/5
BenchmarkDataSync/StateMutex/multi-thread/Mixed-5/5-12                2815182               436.3 ns/op           459 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/Mixed-7/3
BenchmarkDataSync/StateMutex/multi-thread/Mixed-7/3-12                2136508               548.0 ns/op           605 B/op          0 allocs/op
BenchmarkDataSync/StateMutex/multi-thread/Mixed-9/1
BenchmarkDataSync/StateMutex/multi-thread/Mixed-9/1-12                1612203               758.3 ns/op           801 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed
BenchmarkDataSync/StateAtomicFixed/single-thread
BenchmarkDataSync/StateAtomicFixed/single-thread/WriteState
BenchmarkDataSync/StateAtomicFixed/single-thread/WriteState-12             17719366                86.75 ns/op          139 B/op          1 allocs/op
BenchmarkDataSync/StateAtomicFixed/single-thread/ReadState
BenchmarkDataSync/StateAtomicFixed/single-thread/ReadState-12              589778932                2.045 ns/op           0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread
BenchmarkDataSync/StateAtomicFixed/multi-thread/WriteState
BenchmarkDataSync/StateAtomicFixed/multi-thread/WriteState-12               6581025               159.2 ns/op           749 B/op          1 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/ReadState
BenchmarkDataSync/StateAtomicFixed/multi-thread/ReadState-12               1000000000               0.5504 ns/op          0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-1/9
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-1/9-12                4351668               257.0 ns/op          1140 B/op          1 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-3/7
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-3/7-12                2459662               562.3 ns/op          2383 B/op          3 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-5/5
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-5/5-12                1685342               838.7 ns/op          3741 B/op          5 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-7/3
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-7/3-12                1326567              1540 ns/op            7381 B/op          7 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-9/1
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-9/1-12                1000000              1726 ns/op            7824 B/op          9 allocs/op
BenchmarkDataSync/StateAtomic
BenchmarkDataSync/StateAtomic/single-thread
BenchmarkDataSync/StateAtomic/single-thread/WriteState
BenchmarkDataSync/StateAtomic/single-thread/WriteState-12           18016168                86.11 ns/op          137 B/op          1 allocs/op
BenchmarkDataSync/StateAtomic/single-thread/ReadState
BenchmarkDataSync/StateAtomic/single-thread/ReadState-12            579887203                2.090 ns/op           0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread
BenchmarkDataSync/StateAtomic/multi-thread/WriteState
BenchmarkDataSync/StateAtomic/multi-thread/WriteState-12             2889266               447.0 ns/op          1464 B/op          7 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/ReadState
BenchmarkDataSync/StateAtomic/multi-thread/ReadState-12             1000000000               0.5686 ns/op          0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-1/9
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-1/9-12              2639113               454.3 ns/op          1274 B/op          6 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-3/7
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-3/7-12              1000000              1317 ns/op            4210 B/op         20 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-5/5
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-5/5-12               690360              2249 ns/op            7567 B/op         35 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-7/3
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-7/3-12               544074              3003 ns/op            9836 B/op         49 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-9/1
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-9/1-12               384480              3992 ns/op           13659 B/op         64 allocs/op
PASS
ok      github.com/go-auxiliaries/benchs/async/state-consistency        54.004s
?       github.com/go-auxiliaries/benchs/async/state-consistency/code   [no test files]
```

## Summary

### Read flow

You can see that mutex implementation takes `6 times` slower in single-thread read workload:
```
BenchmarkDataSync/StateMutex/single-thread/ReadState-12             98048242                12.19 ns/op            0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/single-thread/ReadState-12       589778932                2.045 ns/op           0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomic/single-thread/ReadState-12            579887203                2.090 ns/op           0 B/op          0 allocs/op
```

For multi-thread read workload `atomic` implementations are at least `90 times` faster.
```
BenchmarkDataSync/StateMutex/multi-thread/ReadState-12              24539293                48.51 ns/op            0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/ReadState-12        1000000000               0.5504 ns/op          0 B/op          0 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/ReadState-12             1000000000               0.5686 ns/op          0 B/op          0 allocs/op
```

### Write flow

For single-thread write workload `atomic` implementations are at least `1.5 times` slower.
```
BenchmarkDataSync/StateMutex/single-thread/WriteState-12            21621788                50.89 ns/op           93 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/single-thread/WriteState-12      17719366                86.75 ns/op          139 B/op          1 allocs/op
BenchmarkDataSync/StateAtomic/single-thread/WriteState-12           18016168                86.11 ns/op          137 B/op          1 allocs/op
```

For multi-thread write workload `atomic` implementations are at least `3 times` slower.
As expected `StateAtomicFixed` shows better performance over atomic implementations, it is slower than `mutex` only `1.5 times`.
```
BenchmarkDataSync/StateMutex/multi-thread/WriteState-12              13214971              102.9 ns/op            97 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/WriteState-12        6581025               159.2 ns/op           749 B/op          1 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/WriteState-12             2889266               447.0 ns/op          1464 B/op          7 allocs/op
```

Worth to mention that due to the nature of `atomic` implementation it creates state allocation on every write attempt operation, 
which is not the case for `mutex` implementation, because it writes into the state.
As expected, the best performance over `atomic` implementations is shown by `StateAtomicFixed`, which uses `mutex` to reduce state submitting collisions, and as result wasted allocations.

### Mixed flow

Obvious disbalance in read and write performance makes it necessary to find a `rate` between read and write, so that you could pick between `atomic` and `mutex` implementations.
Up until write/read rate `1:20` `StateAtomicFixed` is the best choice.
```
BenchmarkDataSync/StateMutex/multi-thread/Mixed-1/20-12              3392697               360.5 ns/op            99 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-1/20-12        3970881               299.3 ns/op           928 B/op          1 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-1/20-12             2563585               475.8 ns/op          1265 B/op          5 allocs/op
```

Having more writes makes `mutex` implementations more preferable.
```
BenchmarkDataSync/StateMutex/multi-thread/Mixed-1/9-12               6304790               188.0 ns/op            83 B/op          0 allocs/op
BenchmarkDataSync/StateAtomicFixed/multi-thread/Mixed-1/9-12         4351668               257.0 ns/op          1140 B/op          1 allocs/op
BenchmarkDataSync/StateAtomic/multi-thread/Mixed-1/9-12              2639113               454.3 ns/op          1274 B/op          6 allocs/op
```

## Conclusion

Over all `atomic` implementations `StateAtomicFixed` is showed best results, without any downsides.
Mutex implementation is preferable for write-dominated workloads, starting from around `1:20` write/read rate.
