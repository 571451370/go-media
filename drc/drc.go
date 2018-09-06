package drc

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type Context struct {
	framebuffer *image.RGBA
	transforms  []f64.Mat4
	styles      []style
	zbuffer     []float64
}

type style struct {
	background color.RGBA
	linewidth  float64
	pointsize  float64
	stroke     color.RGBA
	fill       color.RGBA
	noFill     bool
	noStroke   bool
}

func New(framebuffer *image.RGBA) *Context {
	var M f64.Mat4
	M.Identity()

	r := framebuffer.Bounds()

	c := &Context{}
	c.framebuffer = framebuffer
	c.transforms = append(c.transforms, M)
	c.styles = append(c.styles, c.defaultStyle())
	c.zbuffer = make([]float64, r.Dx()*r.Dy())
	c.Clear()
	return c
}

func (c *Context) defaultStyle() style {
	return style{
		linewidth: 1,
		stroke:    color.RGBA{255, 255, 255, 255},
		fill:      color.RGBA{0, 0, 0, 255},
		noFill:    true,
		noStroke:  false,
	}
}

func (c *Context) Background(col color.RGBA) {
	s := &c.styles[len(c.styles)-1]
	s.background = col
}

func (c *Context) SetStroke(col color.RGBA) {
	s := &c.styles[len(c.styles)-1]
	s.stroke = col
	s.noStroke = false
}

func (c *Context) SetFill(col color.RGBA) {
	s := &c.styles[len(c.styles)-1]
	s.fill = col
	s.noFill = false
}

func (c *Context) NoStroke() {
	s := &c.styles[len(c.styles)-1]
	s.noStroke = true
}

func (c *Context) SetLineWidth(lw float64) {
	s := &c.styles[len(c.styles)-1]
	s.linewidth = lw
}

func (c *Context) SetPointSize(sz float64) {
	s := &c.styles[len(c.styles)-1]
	s.pointsize = sz
}

func (c *Context) LineWidth() float64 {
	s := &c.styles[len(c.styles)-1]
	return s.linewidth
}

func (c *Context) pixel(x, y int, z float64, col color.RGBA) {
	fb := c.framebuffer
	r := fb.Bounds()
	n := y*r.Dx() + x
	if n < len(c.zbuffer) && z <= c.zbuffer[n] {
		c.zbuffer[n] = z
		fb.Set(x, y, col)
	}
}

func (c *Context) pixelRegion(x, y int, z float64, col color.RGBA) {
	s := &c.styles[len(c.styles)-1]
	r := int(s.pointsize / 2)
	for i := -r; i <= r; i++ {
		for j := -r; j <= r; j++ {
			px := x + j
			py := y + i
			c.pixel(px, py, z, col)
		}
	}
}

func (c *Context) Point(x, y int) {
	c.Point3(float64(x), float64(y), 1)
}

func (c *Context) Point3(x, y, z float64) {
	s := &c.styles[len(c.styles)-1]
	c.pixelRegion(int(x+0.5), int(y+0.5), z, s.stroke)
}

func (c *Context) Line(x0, y0, x1, y1 int) {
	c.lineConstantZ(x0, y0, x1, y1, 1)
}

func (c *Context) lineConstantZ(x0, y0, x1, y1 int, z float64) {
	dx := c.abs(x1 - x0)
	dy := c.abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy
	ed := 1.0
	if dx+dy != 0 {
		ed = math.Sqrt(float64(dx*dx) + float64(dy*dy))
	}

	s := &c.styles[len(c.styles)-1]
	wd := (s.linewidth + 1) / 2
	for {
		c.Point3(float64(x0), float64(y0), z)
		e2 := err
		x2 := x0
		if 2*e2 >= -dx {
			e2 += dy
			y2 := y0
			for float64(e2) < ed*wd && (y1 != y2 || dx > dy) {
				y2 += sy
				c.Point3(float64(x0), float64(y2), z)
				e2 += dx
			}
			if x0 == x1 {
				break
			}
			e2 = err
			err -= dy
			x0 += sx
		}
		if 2*e2 <= dy {
			e2 = dx - e2
			for float64(e2) < ed*wd && (x1 != x2 || dx < dy) {
				x2 += sx
				c.Point3(float64(x2), float64(y0), z)
				e2 += dy
			}
			if y0 == y1 {
				break
			}
			err += dx
			y0 += sy
		}
	}
}

func (c *Context) Line3(x0, y0, z0, x1, y1, z1 float64) {
	p0 := f64.Vec3{x0, y0, z0}
	p1 := f64.Vec3{x1, y1, z1}

	m := c.transforms[len(c.transforms)-1]
	p0 = m.Transform3(p0)
	p1 = m.Transform3(p1)

	x0, y0, z0 = p0.X, p0.Y, p0.Z
	x1, y1, z1 = p1.X, p1.Y, p1.Z

	dx := math.Abs(x1 - x0)
	dy := math.Abs(y1 - y0)
	dz := math.Abs(z1 - z0)
	sx, sy, sz := -1.0, -1.0, -1.0
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	if z0 < z1 {
		sz = 1
	}

	dm := math.Max(dx, math.Max(dy, dz))
	i := dm

	x1 = dm / 2
	y1 = x1
	z1 = x1

	lw := c.LineWidth()
	for {
		c.lineConstantZ(int(x0), int(y0), int(x0+lw), int(y0+lw), z0)
		if i--; i <= 0 {
			break
		}
		x1 -= dx
		if x1 < 0 {
			x1 += dm
			x0 += sx
		}
		y1 -= dy
		if y1 < 0 {
			y1 += dm
			y0 += sy
		}
		z1 -= dz
		if z1 < 0 {
			z1 += dm
			z0 += sz
		}
	}
}

func (c *Context) Circle(xm, ym, r int) {
	s := &c.styles[len(c.styles)-1]

loop:
	for n := 0; n < 2; n++ {
		x := -r
		y := 0
		err := 2 - 2*r
		for {
			if n == 0 {
				if s.noFill {
					continue loop
				}
				xs := c.min(xm-x, xm+x)
				xe := c.max(xm-x, xm+x)
				ys := c.min(ym-y, ym+y)
				ye := c.max(ym-y, ym+y)
				for i := xs; i <= xe; i++ {
					for j := ys; j <= ye; j++ {
						c.pixel(i, j, 1, s.fill)
					}
				}
			} else if n == 1 {
				if s.noStroke {
					continue loop
				}
				c.pixelRegion(xm-x, ym+y, 1, s.stroke)
				c.pixelRegion(xm-y, ym-x, 1, s.stroke)
				c.pixelRegion(xm+x, ym-y, 1, s.stroke)
				c.pixelRegion(xm+y, ym+x, 1, s.stroke)
			}

			r := err
			if r <= y {
				y++
				err += y*2 + 1
			}
			if r > x || err > y {
				x++
				err += x*2 + 1
			}

			if x >= 0 {
				break
			}
		}
	}
}

func (c *Context) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *Context) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *Context) abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *Context) PushMatrix(m f64.Mat4) {
	c.transforms = append(c.transforms, m)
}

func (c *Context) PopMatrix() f64.Mat4 {
	n := len(c.transforms) - 1
	m := c.transforms[n]
	c.transforms = c.transforms[:n]
	return m
}

func (c *Context) SetMatrix(m f64.Mat4) {
	p := &c.transforms[len(c.transforms)-1]
	*p = m
}

func (c *Context) Clear() {
	for i := range c.zbuffer {
		c.zbuffer[i] = math.MaxFloat32
	}

	s := &c.styles[len(c.styles)-1]
	fb := c.framebuffer
	draw.Draw(fb, fb.Bounds(), image.NewUniform(s.background), image.ZP, draw.Over)
}
