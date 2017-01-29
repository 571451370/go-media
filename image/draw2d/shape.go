package draw2d

import (
	"image/color"
	"image/draw"
	"math"
)

func Line(img draw.Image, x0, y0, x1, y1 int, color color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	var e2 int
	for {
		img.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
}

func Circle(img draw.Image, cx, cy, radius int, color color.Color) {
	x := radius
	y := 0
	e := 0
	for x >= y {
		img.Set(cx+x, cy+y, color)
		img.Set(cx+y, cy+x, color)
		img.Set(cx-y, cy+x, color)
		img.Set(cx-x, cy+y, color)
		img.Set(cx-x, cy-y, color)
		img.Set(cx-y, cy-x, color)
		img.Set(cx+y, cy-x, color)
		img.Set(cx+x, cy-y, color)

		if e <= 0 {
			y += 1
			e += 2*y + 1
		}
		if e > 0 {
			x -= 1
			e -= 2*x + 1
		}
	}
}

func FilledCircle(img draw.Image, cx, cy, radius int, color color.Color) {
	r2 := radius * radius
	for x := radius; x >= 0; x-- {
		y := int(math.Sqrt(float64(r2-x*x)) + 0.5)
		dx := cx - x
		Line(img, dx, cy-y, dx, cy+y, color)
		dx = cx + x
		Line(img, dx, cy-y, dx, cy+y, color)
	}
}
