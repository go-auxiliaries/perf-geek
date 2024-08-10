package code

type CountGeneric[T CBInterface] struct {
	CB CB
}

func (i *CountGeneric[T]) Count(val string) int {
	out := 0
	for _, ch := range val {
		out += i.CB(ch)
	}
	return out
}
