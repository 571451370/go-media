package mathutil

func GCD(a, b int) int {
	k := Max(a, b)
	m := Min(a, b)
	for m != 0 {
		r := k % m
		k = m
		m = r
	}
	return k
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

func Clamp(x, a, b int) int {
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return x
}