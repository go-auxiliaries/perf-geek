package code

type State struct {
	fruits         []string
	numberOfApples int
	_              int // fix for https://github.com/golang/go/issues/67764, adding this int removes perf degradation
}
