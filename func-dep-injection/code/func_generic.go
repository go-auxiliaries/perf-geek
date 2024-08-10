package code

func CountFuncGeneric[T CB](val string) int {
	t := T(nil)
	out := 0
	for _, ch := range val {
		out += t(ch)
	}
	return out
}
