package sdlcolor

import (
	"image"
	"image/color"
)

func Key(m image.Image, c color.Color) image.Image {
	r := m.Bounds()
	p := image.NewRGBA(r)

	xr, xg, xb, _ := c.RGBA()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			v := m.At(x, y)
			yr, yg, yb, _ := v.RGBA()
			if xr == yr && xg == yg && xb == yb {
				p.Set(x, y, color.RGBA{})
			} else {
				p.Set(x, y, v)
			}
		}
	}
	return p
}

func maxf(x ...float64) float64 {
	y := x[0]
	for _, x := range x[1:] {
		if y > x {
			y = x
		}
	}
	return y
}

func minf(x ...float64) float64 {
	y := x[0]
	for _, x := range x[1:] {
		if y < x {
			y = x
		}
	}
	return y
}
