package imageutil

import (
	"image/color"
	"image/draw"
)

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func Line(m draw.Image, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	e := dx + dy
	for {
		m.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			x0 += sx
		}
		if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}

func Circle(m draw.Image, xc, yc, r int, c color.Color) {
	x, y := r-1, 0
	dx, dy := 1, 1
	e := dx - 2*r
	for x >= y {
		m.Set(xc+x, yc+y, c)
		m.Set(xc+y, yc+x, c)
		m.Set(xc-y, yc+x, c)
		m.Set(xc-x, yc+y, c)
		m.Set(xc-x, yc-y, c)
		m.Set(xc-y, yc-x, c)
		m.Set(xc+y, yc-x, c)
		m.Set(xc+x, yc-y, c)

		if e <= 0 {
			y++
			e += dy
			dy += 2
		}

		if e > 0 {
			x--
			dx += 2
			e += dx - 2*r
		}
	}
}

func Triangle(m draw.Image, x0, y0, x1, y1, x2, y2 int, c color.Color) {
	Line(m, x0, y0, x1, y1, c)
	Line(m, x0, y0, x2, y2, c)
	Line(m, x1, y1, x2, y2, c)
}
