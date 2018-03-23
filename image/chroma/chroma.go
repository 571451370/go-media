package chroma

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/math/f64"
)

type HSL struct {
	H, S, L float64
}

type HSV struct {
	H, S, V float64
}

var (
	HSVModel   = color.ModelFunc(hsvModel)
	HSLModel   = color.ModelFunc(hslModel)
	Vec3dModel = color.ModelFunc(vec3dModel)
	Vec4dModel = color.ModelFunc(vec4dModel)
)

func vec3dModel(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	return f64.Vec3{
		float64(n.R) / 255,
		float64(n.G) / 255,
		float64(n.B) / 255,
	}
}

func vec4dModel(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	return f64.Vec4{
		float64(n.R) / 255,
		float64(n.G) / 255,
		float64(n.B) / 255,
		float64(n.A) / 255,
	}
}

func hsvModel(c color.Color) color.Color {
	b := color.RGBAModel.Convert(c).(color.RGBA)
	return RGB2HSV(b)
}

func hslModel(c color.Color) color.Color {
	b := hsvModel(c).(HSV)
	return HSV2HSL(b)
}

func (h HSV) RGBA() (r, g, b, a uint32) {
	c := HSV2RGB(h)
	return color.RGBA{c.R, c.G, c.B, c.A}.RGBA()
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	c := HSL2HSV(h)
	return c.RGBA()
}

func HSV2RGB(c HSV) color.RGBA {
	h := math.Mod(c.H, 360)
	s := f64.Clamp(c.S, 0, 1)
	v := f64.Clamp(c.V, 0, 1)

	var r color.RGBA
	if s == 0 {
		x := uint8(f64.Clamp(v*255, 0, 255))
		r = color.RGBA{x, x, x, 255}
	} else {
		b := (1 - s) * v
		vb := v - b
		hm := math.Mod(h, 60)

		var cr, cg, cb float64
		switch int(h / 60) {
		case 0:
			cr = v
			cg = vb*h/60 + b
			cb = b
		case 1:
			cr = vb*(60-hm)/60 + b
			cg = v
			cb = b
		case 2:
			cr = b
			cg = v
			cb = vb*hm/60 + b
		case 3:
			cr = b
			cg = vb*(60-hm)/60 + b
			cb = v
		case 4:
			cr = vb*hm/60 + b
			cg = b
			cb = v
		case 5:
			cr = v
			cg = b
			cb = vb*(60-hm)/60 + b
		}
		cr = f64.Clamp(cr*255, 0, 255)
		cg = f64.Clamp(cg*255, 0, 255)
		cb = f64.Clamp(cb*255, 0, 255)
		r = color.RGBA{uint8(cr), uint8(cg), uint8(cb), 255}
	}
	return r
}

func RGB2HSV(c color.RGBA) HSV {
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	min := math.Min(r, math.Max(g, b))
	max := math.Max(r, math.Max(g, b))
	delta := max - min

	v := float64(max) / 255.0
	if delta == 0 {
		return HSV{0, 0, v}
	}

	s := float64(delta) / float64(max)

	h := 0.0
	if r == max {
		h = float64(g-b) / float64(delta)
	} else if g == max {
		h = 2 + float64(b-r)/float64(delta)
	} else {
		h = 4 + float64(r-g)/float64(delta)
	}

	h *= 60
	if h < 0 {
		h += 360
	}
	return HSV{h, s, v}
}

func HSV2HSL(c HSV) HSL {
	h := c.H
	l := (2 - c.S) * c.V
	s := c.S * c.V
	if l <= 1 {
		s /= l
	} else {
		s /= 2 - l
	}
	l /= 2
	return HSL{h, s, l}
}

func HSL2HSV(c HSL) HSV {
	h := c.H
	l := c.L * 2
	s := c.S
	if l <= 1 {
		s *= l
	} else {
		s *= 2 - l
	}
	v := (l + s) / 2
	s = 2 * s / (l + s)
	return HSV{h, s, v}
}

func VEC42RGBA(c f64.Vec4) color.RGBA {
	const eps = 1.001
	if c.X <= eps {
		c.X *= 255
	}
	if c.Y <= eps {
		c.Y *= 255
	}
	if c.Z <= eps {
		c.Z *= 255
	}
	if c.W <= eps {
		c.W *= 255
	}
	c.X = f64.Clamp(c.X, 0, 255)
	c.Y = f64.Clamp(c.Y, 0, 255)
	c.Z = f64.Clamp(c.Z, 0, 255)
	c.W = f64.Clamp(c.W, 0, 255)
	return color.RGBA{
		uint8(c.X),
		uint8(c.Y),
		uint8(c.Z),
		uint8(c.W),
	}
}

func RGBA2VEC4(c color.RGBA) f64.Vec4 {
	return f64.Vec4{
		float64(c.R) / 255.0,
		float64(c.G) / 255.0,
		float64(c.B) / 255.0,
		float64(c.A) / 255.0,
	}
}

func ParseRGBA(s string) (color.RGBA, error) {
	var r, g, b, a uint8
	n, _ := fmt.Sscanf(s, "rgb(%v,%v,%v)", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "rgba(%v,%v,%v,%v)", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x", &r)
	if n == 1 {
		return color.RGBA{r, r, r, 255}, nil
	}

	var h HSV
	n, _ = fmt.Sscanf(s, "hsv(%v,%v,%v)", &h.H, &h.S, &h.V)
	if n == 3 {
		return HSV2RGB(h), nil
	}

	return color.RGBA{}, fmt.Errorf("failed to parse color %q, unknown format", s)
}

func RandRGB() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

func RandRGBA() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
	}
}

func RandHSV() HSV {
	return HSV{
		H: rand.Float64() * 360,
		S: rand.Float64(),
		V: rand.Float64(),
	}
}

func MixRGBA(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		uint8(float64(a.R)*(1-t) + t*float64(b.R)),
		uint8(float64(a.G)*(1-t) + t*float64(b.G)),
		uint8(float64(a.B)*(1-t) + t*float64(b.B)),
		uint8(float64(a.A)*(1-t) + t*float64(b.A)),
	}
}

func MixHSL(a, b HSL, t float64) HSL {
	return HSL{
		a.H*(1-t) + t*b.H,
		a.S*(1-t) + t*b.S,
		a.L*(1-t) + t*b.L,
	}
}

func RGBA32(c color.RGBA) uint32 {
	return uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(c.A)<<24
}

func BGRA32(c color.RGBA) uint32 {
	return uint32(c.B) | uint32(c.G)<<8 | uint32(c.R)<<16 | uint32(c.A)<<24
}
