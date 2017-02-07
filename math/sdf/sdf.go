// http://iquilezles.org/www/articles/distfunctions/distfunctions.htm

package sdf

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func Sphere(p f64.Vec3, s float64) float64 {
	return p.Len() - s
}

func Box(p, b f64.Vec3) float64 {
	d := p.Abs().Sub(b)
	return math.Min(d.MaxComp(), 0) + d.Max(0).Len()
}

func Torus(p f64.Vec3, t f64.Vec2) float64 {
	l := f64.Vec2{p.X, p.Z}.Len()
	q := f64.Vec2{l - t.X, p.Y}
	return q.Len() - t.Y
}

func Cylinder(p f64.Vec3, h f64.Vec2) float64 {
	l := f64.Vec2{p.X, p.Z}.Len()
	q := f64.Vec2{l, p.Y}.Abs()
	d := q.Sub(h)
	t1 := math.Min(d.MaxComp(), 0)
	t2 := d.Max(0).Len()
	return t1 + t2
}

func Cone(p, c f64.Vec3) float64 {
	l := f64.Vec2{p.X, p.Z}.Len()
	q := f64.Vec2{l, p.Y}
	d1 := -q.Y - c.Z
	d2 := math.Max(q.Dot(f64.Vec2{c.X, c.Y}), q.Y)
	t1 := f64.Vec2{d1, d2}.Max(0).Len()
	t2 := math.Min(math.Max(d1, d2), 0)
	return t1 + t2
}

func Plane(p f64.Vec3, n f64.Vec4) float64 {
	return p.Dot(n.V3()) + n.W
}

func OpS1(d1, d2 float64) float64 {
	return math.Max(-d2, d1)
}

func OpU1(d1, d2 float64) float64 {
	return math.Min(d1, d2)
}

func OpI1(d1, d2 float64) float64 {
	return math.Max(d1, d2)
}

func OpU2(d1, d2 f64.Vec2) f64.Vec2 {
	if d1.X < d2.X {
		return d1
	}
	return d2
}
