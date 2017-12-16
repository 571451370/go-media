package chroma

import (
	"fmt"
	"image/color"
	"math"

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
	return rgb2hsv(b)
}

func hslModel(c color.Color) color.Color {
	b := hsvModel(c).(HSV)
	return hsv2hsl(b)
}

func (h HSV) RGBA() (r, g, b, a uint32) {
	c := hsv2rgb(h)
	return color.RGBA{c.R, c.G, c.B, c.A}.RGBA()
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	c := hsl2hsv(h)
	return c.RGBA()
}

func hsv2rgb(c HSV) color.RGBA {
	h := math.Mod(c.H, 360)
	s := clampf(c.S, 0, 1)
	v := clampf(c.V, 0, 1)

	var r color.RGBA
	if s == 0 {
		x := uint8(clampf(v*255, 0, 255))
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
		cr = clampf(cr*255, 0, 255)
		cg = clampf(cg*255, 0, 255)
		cb = clampf(cb*255, 0, 255)
		r = color.RGBA{uint8(cr), uint8(cg), uint8(cb), 255}
	}
	return r
}

func rgb2hsv(c color.RGBA) HSV {
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	min := minf(r, minf(g, b))
	max := maxf(r, maxf(g, b))
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

func hsv2hsl(c HSV) HSL {
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

func hsl2hsv(c HSL) HSV {
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

func min8(x, y uint8) uint8 {
	if x < y {
		return x
	}
	return y
}

func minf(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func max8(x, y uint8) uint8 {
	if x > y {
		return x
	}
	return y
}

func maxf(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func clampf(x, a, b float64) float64 {
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return x
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

	return color.RGBA{}, fmt.Errorf("failed to parse color %q, unknown format", s)
}
