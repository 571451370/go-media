package sdlcolor

import (
	"github.com/qeedquan/go-gfx/math/f64"
	"github.com/qeedquan/go-gfx/sdl"
)

var (
	White       = sdl.Color{255, 255, 255, 255}
	Black       = sdl.Color{0, 0, 0, 255}
	Transparent = sdl.Color{0, 0, 0, 0}
)

type HSL struct {
	H, S, L float64
}

type HSV struct {
	H, S, V float64
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
