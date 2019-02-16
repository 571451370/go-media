package imageutil

import (
	"image/color"
	"image/draw"
	"math"
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

func ThickLine(m draw.Image, x0, y0, x1, y1 int, wd float64, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	x2, y2 := 0, 0
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	e := dx - dy
	e2 := 0
	ed := 1.0
	if dx+dy != 0 {
		ed = math.Hypot(float64(dx), float64(dy))
	}

	cc := color.RGBAModel.Convert(c).(color.RGBA)
	wd = (wd + 1) / 2
	for {
		r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		col := color.RGBA{r, g, b, 255}
		m.Set(x0, y0, col)

		e2 = e
		x2 = x0
		if 2*e2 >= -dx {
			e2 += dy
			y2 = y0
			for float64(e2) < ed*wd && (y1 != y2 || dx > dy) {
				y2 += sy
				r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e2))/ed-wd+1)))
				g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e2))/ed-wd+1)))
				b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e2))/ed-wd+1)))
				col := color.RGBA{r, g, b, 255}
				m.Set(x0, y2, col)
				e2 += dx
			}
			if x0 == x1 {
				break
			}
			e2 = e
			e -= dy
			x0 += sx
		}
		if 2*e2 <= dy {
			e2 = dx - e2
			for float64(e2) < ed*wd && (x1 != x2 || dx < dy) {
				x2 += sx
				r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e2))/ed-wd+1)))
				g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e2))/ed-wd+1)))
				b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e2))/ed-wd+1)))
				col := color.RGBA{r, g, b, 255}
				m.Set(x2, y0, col)
				e2 += dy
			}
			if y0 == y1 {
				break
			}
			e += dx
			y0 += sy
		}
	}
}

func Circle(m draw.Image, cx, cy, r int, c color.Color) {
	x := -r
	y := 0
	e := 2 - 2*r
	for {
		m.Set(cx-x, cy+y, c)
		m.Set(cx-y, cy-x, c)
		m.Set(cx+x, cy-y, c)
		m.Set(cx+y, cy+x, c)
		r := e
		if r <= y {
			y++
			e += 2*y + 1
		}
		if r > x || e > y {
			x++
			e += 2*x + 1
		}
		if x > 0 {
			break
		}
	}
}

func FilledCircle(m draw.Image, cx, cy, r int, c color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				m.Set(cx+x, cy+y, c)
			}
		}
	}
}

func Triangle(m draw.Image, x0, y0, x1, y1, x2, y2 int, c color.Color) {
	Line(m, x0, y0, x1, y1, c)
	Line(m, x0, y0, x2, y2, c)
	Line(m, x1, y1, x2, y2, c)
}

func FilledTriangle(m draw.Image, x0, y0, x1, y1, x2, y2 int, c color.Color) {
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y0 > y2 {
		x0, y0, x2, y2 = x2, y2, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}
	filledTriangleBottom(m, x0, y0, x1, y1, x2, y2, c)
	filledTriangleTop(m, x0, y0, x1, y1, x2, y2, c)
}

func filledTriangleBottom(m draw.Image, x0, y0, x1, y1, x2, y2 int, c color.Color) {
	i1 := float64(x1-x0) / float64(y1-y0)
	i2 := float64(x2-x0) / float64(y2-y0)

	cx1 := float64(x0)
	cx2 := cx1
	for y := y0; y <= y1; y++ {
		Line(m, int(cx1), y, int(cx2), y, c)
		cx1 += i1
		cx2 += i2
	}
}

func filledTriangleTop(m draw.Image, x0, y0, x1, y1, x2, y2 int, c color.Color) {
	i1 := float64(x2-x0) / float64(y2-y0)
	i2 := float64(x2-x1) / float64(y2-y1)

	cx1 := float64(x2)
	cx2 := cx1
	for y := y2; y > y0; y-- {
		Line(m, int(cx1), y, int(cx2), y, c)
		cx1 -= i1
		cx2 -= i2
	}
}

func EllipseRect(m draw.Image, x0, y0, x1, y1 int, c color.Color) {
	a := abs(x1 - x0)
	b := abs(y1 - y0)
	b1 := b & 1
	dx := 4 * (1 - a) * b * b
	dy := 4 * (b1 + 1) * a * a
	e := dx + dy + b1*a*a
	if x0 > x1 {
		x0 = x1
		x1 += a
	}
	if y0 > y1 {
		y0 = y1
	}
	y0 += (b + 1) / 2
	y1 = y0 - b1
	a *= 8 * a
	b1 = 8 * b * b
	for {
		m.Set(x1, y0, c)
		m.Set(x0, y0, c)
		m.Set(x0, y1, c)
		m.Set(x1, y1, c)
		e2 := 2 * e
		if e2 <= dy {
			y0++
			y1--
			dy += a
			e += dy
		}
		if e2 >= dx || 2*e > dy {
			x0++
			x1--
			dx += b1
			e += dx
		}
		if x0 > x1 {
			break
		}
	}
	for y0-y1 < b {
		m.Set(x0-1, y0, c)
		m.Set(x1+1, y0, c)
		y0++
		m.Set(x0-1, y1, c)
		m.Set(x1+1, y1, c)
		y1--
	}
}
