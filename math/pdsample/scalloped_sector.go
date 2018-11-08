package pdsample

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/math/f64"
)

type ArcData struct {
	P               f64.Vec2
	R               float64
	Sign            float64
	D               float64
	Theta           float64
	IntegralAtStart float64
	RSquared        float64
	DSquared        float64
}

type ScallopedSector struct {
	P            f64.Vec2
	A1, A2, Area float64
	Arcs         [2]ArcData
}

type ScallopedRegion struct {
	regions []ScallopedSector
	minArea float64
	area    float64
}

func NewScallopedSector(pt f64.Vec2, a1, a2 float64, p1 f64.Vec2, r1, sign1 float64, p2 f64.Vec2, r2, sign2 float64) *ScallopedSector {
	v1 := p1.Sub(pt)
	v2 := p2.Sub(pt)

	arc1 := ArcData{
		P:     p1,
		R:     r1,
		Sign:  sign1,
		D:     math.Sqrt(v1.X*v1.X + v1.Y*v1.Y),
		Theta: math.Atan2(v1.Y, v1.X),
	}
	arc1.RSquared = arc1.R * arc1.R
	arc1.DSquared = arc1.D * arc1.D
	arc1.IntegralAtStart = integralOfDistToCircle(a1-arc1.Theta, arc1.D, arc1.R, arc1.Sign)

	arc2 := ArcData{
		P:     p2,
		R:     r2,
		Sign:  sign2,
		D:     math.Sqrt(v2.X*v2.X + v2.Y*v2.Y),
		Theta: math.Atan2(v2.Y, v2.X),
	}
	arc2.RSquared = arc2.R * arc2.R
	arc2.DSquared = arc2.D * arc2.D
	arc2.IntegralAtStart = integralOfDistToCircle(a1-arc2.Theta, arc2.D, arc2.R, arc2.Sign)

	s := &ScallopedSector{
		P:    pt,
		A1:   a1,
		A2:   a2,
		Arcs: [2]ArcData{arc1, arc2},
	}
	s.Area = s.calcAreaToAngle(a2)
	return s
}

func (s *ScallopedSector) calcAreaToAngle(angle float64) float64 {
	underInner := integralOfDistToCircle(angle-s.Arcs[0].Theta, s.Arcs[0].D, s.Arcs[0].R, s.Arcs[0].Sign) - s.Arcs[0].IntegralAtStart
	underOuter := integralOfDistToCircle(angle-s.Arcs[1].Theta, s.Arcs[1].D, s.Arcs[1].R, s.Arcs[1].Sign) - s.Arcs[1].IntegralAtStart
	return underOuter - underInner
}

func (s *ScallopedSector) calcAngleForArea(area float64) float64 {
	lo := s.A1
	hi := s.A2
	cur := lo + (hi-lo)*rand.Float64()

	for i := 0; i < 10; i++ {
		if s.calcAreaToAngle(cur) < s.Area {
			lo = cur
			cur = (cur + hi) * .5
		} else {
			hi = cur
			cur = (lo + cur) * .5
		}
	}
	return cur
}

func (s *ScallopedSector) distToCurve(angle float64, index int) float64 {
	alpha := angle - s.Arcs[index].Theta
	sin_alpha := math.Sin(alpha)
	t0 := s.Arcs[index].RSquared - s.Arcs[index].DSquared*sin_alpha*sin_alpha
	if t0 < 0 {
		return s.Arcs[index].D * math.Cos(alpha)
	}
	return s.Arcs[index].D*math.Cos(alpha) + s.Arcs[index].Sign*math.Sqrt(t0)
}

func (s *ScallopedSector) Sample() f64.Vec2 {
	angle := s.calcAngleForArea(s.Area * rand.Float64())
	d1 := s.distToCurve(angle, 0)
	d2 := s.distToCurve(angle, 1)
	d := math.Sqrt(d1*d1 + (d2*d2-d1*d1)*rand.Float64())
	return f64.Vec2{
		s.P.X + math.Cos(angle)*d,
		s.P.Y + math.Sin(angle)*d,
	}
}

func (s *ScallopedSector) canonizeAngle(angle float64) float64 {
	delta := math.Mod(angle-s.A1, 2*math.Pi)
	if delta < 0 {
		delta += 2 * math.Pi
	}
	return s.A1 + delta
}

func (s *ScallopedSector) distToCircle(angle float64, C f64.Vec2, r float64) (d1, d2 float64) {
	v := f64.Vec2{C.X - s.P.X, C.Y - s.P.Y}
	dSquared := v.X*v.X + v.Y*v.Y
	theta := math.Atan2(v.Y, v.X)
	alpha := angle - theta
	sin_alpha := math.Sin(alpha)
	xSquared := r*r - dSquared*sin_alpha*sin_alpha

	if xSquared < 0 {
		d1 = -10000000
		d2 = d1
	} else {
		a := math.Sqrt(dSquared) * math.Cos(alpha)
		x := math.Sqrt(xSquared)
		d1 = a - x
		d2 = a + x
	}

	return
}

func (s *ScallopedSector) SubtractDisk(C f64.Vec2, r float64) (regions []ScallopedSector) {
	return
}

func integralOfDistToCircle(x, d, r, k float64) float64 {
	sin_x := math.Sin(x)
	d_sin_x := d * sin_x
	y := sin_x * d / r
	if y < -1 {
		y = -1
	} else if y > 1 {
		y = 1
	}

	theta := math.Asin(y)

	return (r*(r*(x+
		k*theta)+
		k*math.Cos(theta)*d_sin_x) +
		d*math.Cos(x)*d_sin_x) * .5
}