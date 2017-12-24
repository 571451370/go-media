package mathutil

func GCD(a, b int) int {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}

func LCM(a, b int) int {
	return Abs(a*b) / GCD(a, b)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
