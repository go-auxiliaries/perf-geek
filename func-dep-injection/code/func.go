package code

func CountFunc(val string, cb CB) int {
	out := 0
	for _, ch := range val {
		out += cb(ch)
	}
	return out
}
