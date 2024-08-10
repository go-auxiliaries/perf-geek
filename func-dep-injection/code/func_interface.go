package code

func CountFuncInterface(val string, cb CBInterface) int {
	out := 0
	for _, ch := range val {
		out += cb.CB(ch)
	}
	return out
}
