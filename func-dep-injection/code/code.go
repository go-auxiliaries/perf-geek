package code

type CB func(rune) int

type CBInterface interface {
	CB(rune) int
}
