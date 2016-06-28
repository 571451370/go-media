package sdlcolor

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/sdl"
)

var (
	White       = sdl.Color{255, 255, 255, 255}
	Black       = sdl.Color{0, 0, 0, 255}
	Transparent = sdl.Color{0, 0, 0, 0}
	Red         = sdl.Color{255, 0, 0, 255}
	Blue        = sdl.Color{0, 0, 255, 255}
	Green       = sdl.Color{0, 255, 0, 255}
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
	n := color.RGBAModel.Convert(c).(color.RGBA)
	r := float64(n.R) / 255
	g := float64(n.G) / 255
	b := float64(n.B) / 255

	max := maxf(r, g, b)
	min := minf(r, g, b)
	eps := 1e-4
	d := max - min

	var h, s, v float64

	v = max
	if max != 0 {
		s = d / max
	}
	if math.Abs(max-min) < eps {
		h = 0
	} else {
		switch {
		case math.Abs(max-r) < eps:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case math.Abs(max-g) < eps:
			h = (b-r)/d + 2
		case math.Abs(max-b) < eps:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return HSV{h, s, v}
}

func hslModel(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	r := float64(n.R) / 255
	g := float64(n.G) / 255
	b := float64(n.B) / 255

	max := maxf(r, g, b)
	min := minf(r, g, b)
	eps := 1e-4

	var h, s, l float64
	l = (max + min) / 2
	if math.Abs(max-min) < eps {
		h, s, l = 0, 0, 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}
		switch {
		case math.Abs(max-r) < eps:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case math.Abs(max-g) < eps:
			h = (b-r)/d + 2
		case math.Abs(max-b) < eps:
			h = (r-g)/d + 4
		}

		h /= 6
	}

	return HSL{h, s, l}
}

func (h HSV) RGBA() (r, g, b, a uint32) {
	hue := (2 - h.S) * h.V
	if hue >= 1 {
		hue = 2 - hue
	}
	sat := h.S * h.V / hue
	hsl := HSL{h.H, sat, hue / 2}
	return hsl.RGBA()
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	var c f64.Vec3

	if h.S == 0 {
		c.X, c.Y, c.Z = h.L, h.L, h.L
	} else {
		var q float64
		if h.L < 0.5 {
			q = h.L * (1 + h.S)
		} else {
			q = h.L + h.S - h.L*h.S
		}
		p := 2*h.L - q
		c.X = hue2rgb(p, q, h.H+1.0/3)
		c.Y = hue2rgb(p, q, h.H)
		c.Z = hue2rgb(p, q, h.H-1.0/3)
	}

	return c.RGBA()
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
