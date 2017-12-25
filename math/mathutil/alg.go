package mathutil

func GCD(a, b int) int {
	for a != 0 && b != 0 {
		if a > b {
			a %= b
		} else {
			b %= a
		}
	}

	return Max(a, b)
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

func Abs8(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
