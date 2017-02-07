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
	q := f64.Vec2{p.XZ().Len() - t.X, p.Y}
	return q.Len() - t.Y
}

func Cylinder(p, c f64.Vec3) float64 {
	t1 := p.XZ()
	t2 := c.XY()
	t3 := t1.Distance(t2)
	return t3 - c.Z
}

func Cone(p f64.Vec3, c f64.Vec2) float64 {
	q := p.XY().Len()
	return c.Dot(f64.Vec2{q, p.Z})
}

func HexPrism(p f64.Vec3, h f64.Vec2) float64 {
	q := p.Abs()
	return math.Max(q.Z-h.Y, math.Max((q.X*0.866025+q.Y*0.5), q.Y)-h.X)
}

func Plane(p f64.Vec3, n f64.Vec4) float64 {
	return p.Dot(n.XYZ()) + n.W
}

func Triangle(p, a, b, c f64.Vec3) float64 {
	ba := b.Sub(a)
	pa := p.Sub(a)
	cb := c.Sub(b)
	pb := p.Sub(b)
	ac := a.Sub(c)
	pc := p.Sub(c)
	nor := ba.Cross(ac)

	s1 := f64.Sign(ba.Cross(nor).Dot(pa))
	s2 := f64.Sign(cb.Cross(nor).Dot(pb))
	s3 := f64.Sign(ac.Cross(nor).Dot(pc))
	if math.Sqrt(s1+s2+s3) < 2 {
		t1 := f64.Clamp(ba.Dot(pa)/ba.Dot(ba), 0, 1)
		x1 := ba.Scale(t1).Sub(pa)
		p1 := x1.Dot(x1)

		t2 := f64.Clamp(cb.Dot(pb)/cb.Dot(cb), 0, 1)
		x2 := cb.Scale(t2).Sub(pb)
		p2 := x2.Dot(x2)

		t3 := f64.Clamp(ac.Dot(pc)/ac.Dot(ac), 0, 1)
		x3 := ac.Scale(t3).Sub(pc)
		p3 := x3.Dot(x3)

		return math.Min(p1, math.Min(p2, p3))
	}

	t1 := nor.Dot(pa)
	t2 := nor.Dot(nor)
	return t1 * t1 / t2
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
