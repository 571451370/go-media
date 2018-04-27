package imgui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

func InvLength(lhs f64.Vec2, fail_value float64) float64 {
	d := lhs.X*lhs.X + lhs.Y*lhs.Y
	if d > 0.0 {
		return 1.0 / math.Sqrt(d)
	}
	return fail_value
}

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

func UpperPowerOfTwo(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}

type stbCompress struct {
	barrier, barrier2, barrier3, barrier4 uint32
	din                                   uint32
	dout                                  uint32
}

func (c *stbCompress) in2(i []byte, x uint32) uint32 {
	return (uint32(i[x]) << 8) + uint32(i[(x)+1])
}

func (c *stbCompress) in3(i []byte, x uint32) uint32 {
	return (uint32(i[x]) << 16) + c.in2(i, x+1)
}

func (c *stbCompress) in4(i []byte, x uint32) uint32 {
	return (uint32(i[x]) << 24) + c.in3(i, x+1)
}

func (c *stbCompress) match(output, data []byte, pos, length uint32) {
	// INVERSE of memmove... write each byte before copying the next...
	assert(c.dout+length <= c.barrier)
	if c.dout+length > c.barrier {
		c.dout += length
		return
	}
	if pos < c.barrier4 {
		c.dout = c.barrier + 1
		return
	}
	for i := uint32(0); length != 0; length-- {
		output[c.dout] = data[pos+i]
		c.dout++
		i++
	}
}

func (c *stbCompress) lit(output, data []byte, pos, length uint32) {
	assert(c.dout+length <= c.barrier)
	if c.dout+length > c.barrier {
		c.dout += length
		return
	}
	if pos < c.barrier2 {
		c.dout = c.barrier + 1
		return
	}
	copy(output[c.dout:], data[pos:pos+length])
	c.dout += length
}

func (c *stbCompress) decompressToken(output, input []byte) uint32 {
	i := c.din
	// use fewer if's for cases that expand small
	if input[i] >= 0x20 {
		if input[i] >= 0x80 {
			c.match(output, output, c.dout-uint32(input[i+1])-1, uint32(input[i+0])-0x80+1)
			i += 2
		} else if input[i] >= 0x40 {
			c.match(output, output, c.dout-(c.in2(input[i:], 0)-0x4000+1), uint32(input[i+2]+1))
			i += 3
		} else {
			c.lit(output, input, i+1, uint32(input[i+0]-0x20+1))
			i += 1 + (uint32(input[i]) - 0x20 + 1)
		}
	} else {
		// more ifs for cases that expand large, since overhead is amortized
		if input[i] >= 0x18 {
			c.match(output, output, c.dout-(c.in3(input[i:], 0)-0x180000+1), uint32(input[i+3])+1)
			i += 4
		} else if input[i] >= 0x10 {
			c.match(output, output, c.dout-(c.in3(input[i:], 0)-0x100000+1), c.in2(input[i:], 3)+1)
			i += 5
		} else if input[i] >= 0x08 {
			c.lit(output, input, i+2, c.in2(input[i:], 0)-0x0800+1)
			i += 2 + (c.in2(input[i:], 0) - 0x0800 + 1)
		} else if input[i] == 0x07 {
			c.lit(output, input, i+3, c.in2(input[i:], 1)+1)
			i += 3 + (c.in2(input[i:], 1) + 1)
		} else if input[i] == 0x06 {
			c.match(output, output, c.dout-(c.in3(input[i:], 1)+1), uint32(input[i+4])+1)
			i += 5
		} else if input[i] == 0x04 {
			c.match(output, output, c.dout-(c.in3(input[i:], 1)+1), c.in2(input[i:], 4)+1)
			i += 6
		}
	}
	return i
}

func (c *stbCompress) adler32(adler32 uint32, buffer []byte) uint32 {
	const ADLER_MOD = 65521
	s1 := adler32 & 0xffff
	s2 := adler32 >> 16

	buflen := len(buffer)
	blocklen := buflen % 5552
	var o, i int
	for buflen != 0 {
		for i = 0; i+7 < blocklen; i += 8 {
			for j := 0; j < 8; j++ {
				s1 += uint32(buffer[o+j])
				s2 += s1
			}
			o += 8
		}

		for ; i < blocklen; i++ {
			s1 += uint32(buffer[o])
			s2 += s1
			o++
		}

		s1 %= ADLER_MOD
		s2 %= ADLER_MOD
		buflen -= blocklen
		blocklen = 5552
	}

	return s2<<16 + s1
}

func (c *stbCompress) DecompressLength(input []byte) uint32 {
	return (uint32(input[8]) << 24) + (uint32(input[9]) << 16) + (uint32(input[10]) << 8) + uint32(input[11])
}

func (c *stbCompress) Decompress(output, input []byte) error {
	if c.in4(input, 0) != 0x57bc0000 {
		return fmt.Errorf("invalid header")
	}
	// error! stream is > 4 GB
	if c.in4(input, 4) != 0 {
		return fmt.Errorf("stream too big")
	}

	olen := c.DecompressLength(input)
	c.barrier2 = 0
	c.barrier3 = uint32(len(input))
	c.barrier = olen
	c.barrier4 = 0

	c.din = 16
	c.dout = 0
	for {
		old_i := c.din
		c.din = c.decompressToken(output, input)
		if c.din == old_i {
			if input[c.din] == 0x05 && input[c.din+1] == 0xfa {
				assert(c.dout == olen)
				if c.dout != olen {
					return fmt.Errorf("corruption of olen")
				}
				if c.adler32(1, output[:olen]) != c.in4(input[c.din:], 2) {
					return fmt.Errorf("adler32 mismatch")
				}
				return nil
			} else {
				panic("unreachable")
			}
		}

		if c.dout > olen {
			return nil
		}
	}

	return nil
}

// Parse display precision back from the display format string
func ParseFormatPrecision(format string, default_precision int) int {
	precision := default_precision
	for {
		n := strings.IndexRune(format, '%')
		if n < 0 {
			break
		}
		format = format[n:]

		// Ignore "%%"
		if strings.HasPrefix(format, "%") {
			format = format[1:]
			continue
		}

		for ; len(format) > 0; format = format[1:] {
			if !('0' <= format[0] && format[0] <= '9') {
				break
			}
		}

		if strings.HasPrefix(format, ".") {
			format = format[1:]
			precision, _ = strconv.Atoi(format)
			if precision < 0 || precision > 10 {
				precision = default_precision
			}
		}

		// Maximum precision with scientific notation
		if strings.HasPrefix(format, "e") || strings.HasPrefix(format, "E") {
			precision = -1
		}
		break
	}
	return precision
}

func DataTypeFormatStringCustom(data interface{}, format string) string {
	return fmt.Sprintf(format, data)
}

func DataTypeFormatString(data interface{}, decimal_precision int) string {
	switch v := data.(type) {
	case int:
		if decimal_precision < 0 {
			return fmt.Sprintf("%d", v)
		} else {
			return fmt.Sprintf("%.*d", v)
		}
	case float32:
		// Ideally we'd have a minimum decimal precision of 1 to visually denote that it is a float, while hiding non-significant digits?
		if decimal_precision < 0 {
			return fmt.Sprintf("%f", v)
		} else {
			return fmt.Sprintf("%.*f", v)
		}
	case float64:
		if decimal_precision < 0 {
			return fmt.Sprintf("%f", v)
		} else {
			return fmt.Sprintf("%.*f", v)
		}
	default:
		panic("unreachable")
	}
}

// User can input math operators (e.g. +100) to edit a numerical values.
// NB: This is _not_ a full expression evaluator. We should probably add one though..
func DataTypeApplyOpFromText(buf []byte, initial_value_buf string, data interface{}, scalar_format string) bool {
	return false
}

func TriangleContainsPoint(a, b, c, p f64.Vec2) bool {
	b1 := ((p.X-b.X)*(a.Y-b.Y) - (p.Y-b.Y)*(a.X-b.X)) < 0.0
	b2 := ((p.X-c.X)*(b.Y-c.Y) - (p.Y-c.Y)*(b.X-c.X)) < 0.0
	b3 := ((p.X-a.X)*(c.Y-a.Y) - (p.Y-a.Y)*(c.X-a.X)) < 0.0
	return ((b1 == b2) && (b2 == b3))
}

func truth(x bool) int {
	if x {
		return 1
	}
	return 0
}

func CharIsSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == 0x3000
}

// find beginning-of-line
func StrbolW(buf []rune, buf_mid_line, buf_begin int) int {
	for buf_mid_line > buf_begin && buf[buf_mid_line-1] != '\n' {
		buf_mid_line--
	}
	return buf_mid_line
}

func Acos01(x float64) float64 {
	if x <= 0.0 {
		return math.Pi * 0.5
	}
	if x >= 1.0 {
		return 0.0
	}
	return math.Acos(x)
}

func GetMinimumStepAtDecimalPrecision(decimal_precision int) float64 {
	min_steps := [10]float64{1.0, 0.1, 0.01, 0.001, 0.0001, 0.00001, 0.000001, 0.0000001, 0.00000001, 0.000000001}
	if decimal_precision >= 0 && decimal_precision < 10 {
		return min_steps[decimal_precision]
	}
	return math.Pow(10, float64(-decimal_precision))
}

func F32_TO_INT8_UNBOUND(v float64) int {
	if v >= 0 {
		return int(v*255 + 0.5)
	}
	return int(v*255 - 0.5)
}

func Rotate(v f64.Vec2, cos_a, sin_a float64) f64.Vec2 {
	return f64.Vec2{v.X*cos_a - v.Y*sin_a, v.X*sin_a + v.Y*cos_a}
}

func TriangleBarycentricCoords(a, b, c, p f64.Vec2) (out_u, out_v, out_w float64) {
	v0 := b.Sub(a)
	v1 := c.Sub(a)
	v2 := p.Sub(a)
	denom := v0.X*v1.Y - v1.X*v0.Y
	out_v = (v2.X*v1.Y - v1.X*v2.Y) / denom
	out_w = (v0.X*v2.Y - v2.X*v0.Y) / denom
	out_u = 1.0 - out_v - out_w
	return
}

func TriangleClosestPoint(a, b, c, p f64.Vec2) f64.Vec2 {
	proj_ab := LineClosestPoint(a, b, p)
	proj_bc := LineClosestPoint(b, c, p)
	proj_ca := LineClosestPoint(c, a, p)
	dist2_ab := p.Sub(proj_ab).LenSquared()
	dist2_bc := p.Sub(proj_bc).LenSquared()
	dist2_ca := p.Sub(proj_ca).LenSquared()
	m := math.Min(dist2_ab, math.Min(dist2_bc, dist2_ca))
	if m == dist2_ab {
		return proj_ab
	}
	if m == dist2_bc {
		return proj_bc
	}
	return proj_ca
}

func LineClosestPoint(a, b, p f64.Vec2) f64.Vec2 {
	ap := p.Sub(a)
	ab_dir := b.Sub(a)
	dot := ap.X*ab_dir.X + ap.Y*ab_dir.Y
	if dot < 0.0 {
		return a
	}
	ab_len_sqr := ab_dir.X*ab_dir.X + ab_dir.Y*ab_dir.Y
	if dot > ab_len_sqr {
		return b
	}
	v := ab_dir.Scale(dot / ab_len_sqr)
	v = a.Add(v)
	return v
}