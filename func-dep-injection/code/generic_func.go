package code

type CountGenericFunc[T CBInterface] CB

func (i CountGenericFunc[T]) Count(val string) int {
	out := 0
	for _, ch := range val {
		out += i(ch)
	}
	return out
}
