## Problem

Investigate performance impact of different ways to inject your code.
Consider following example:
```go
func CountSymbols(val string) {
	out := 0
	for _, ch := range val {
	    if ch == 'a' {
			out++
	    }
	}
	return out
}
```

Say you want to generalize it and do something like:
```go
func CountSymbols(val string, cb func(rune) int) {
	out := 0
	for _, ch := range val {
		out += cb(ch)
	    }
	}
	return out
}
```




## Solution
### Interface
```go
type CounterCB interface {
	CB(rune) int
} 

func CountSymbols(val string, cb func(rune) int) {
	out := 0
	for _, ch := range val {
		out += cb(ch)
	    }
	}
	return out
}
```
### func type
```go
func CountSymbols(val string, cb func(rune) int) {
	out := 0
	for _, ch := range val {
		out += cb(ch)
	    }
	}
	return out
}
```

### Generics

## Results
## Summary
## Conclusion

Over all `atomic` implementations `StateAtomicFixed` is showed best results, without any downsides.
Mutex implementation is preferable for write-dominated workloads, starting from around `1:20` write/read rate.
