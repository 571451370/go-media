package imgui

func Decode85Byte(c byte) uint {
	if c >= '\\' {
		return uint(c - 36)
	}
	return uint(c - 35)
}

func Decode85(src, dst []byte) {
	i, j := 0, 0
	for i < len(src) {
		tmp := Decode85Byte(src[i]) + 85*(Decode85Byte(src[i+1])+85*(Decode85Byte(src[i+2])+85*(Decode85Byte(src[i+3])+85*Decode85Byte(src[i+4]))))
		// We can't assume little-endianness
		dst[j] = byte((tmp >> 0) & 0xFF)
		dst[j+1] = byte((tmp >> 8) & 0xFF)
		dst[j+2] = byte((tmp >> 16) & 0xFF)
		dst[j+3] = byte((tmp >> 24) & 0xFF)
		i += 5
		j += 4
	}
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}
