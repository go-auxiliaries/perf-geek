package code

type CountInterface struct {
	CB CB
}

func (i *CountInterface) Count(val string) int {
	out := 0
	for _, ch := range val {
		out += i.CB(ch)
	}
	return out
}
