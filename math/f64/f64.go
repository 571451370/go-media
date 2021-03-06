package f64

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/cmplx"
	"sort"
)

type Vec2 struct {
	X, Y float64
}

func (p Vec2) Add(q Vec2) Vec2 {
	return Vec2{p.X + q.X, p.Y + q.Y}
}

func (p Vec2) Sub(q Vec2) Vec2 {
	return Vec2{p.X - q.X, p.Y - q.Y}
}

func (p Vec2) AddScale(q Vec2, t float64) Vec2 {
	return p.Add(q.Scale(t))
}

func (p Vec2) SubScale(q Vec2, t float64) Vec2 {
	return p.Sub(q.Scale(t))
}

func (p Vec2) Neg() Vec2 {
	return Vec2{-p.X, -p.Y}
}

func (p Vec2) Abs() Vec2 {
	return Vec2{math.Abs(p.X), math.Abs(p.Y)}
}

func (p Vec2) Dot(q Vec2) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Vec2) MinComp() float64 {
	return math.Min(p.X, p.Y)
}

func (p Vec2) MaxComp() float64 {
	return math.Max(p.X, p.Y)
}

func (p Vec2) Max(q Vec2) Vec2 {
	return Vec2{
		math.Max(p.X, q.X),
		math.Max(p.Y, q.Y),
	}
}

func (p Vec2) Min(q Vec2) Vec2 {
	return Vec2{
		math.Min(p.X, q.X),
		math.Min(p.Y, q.Y),
	}
}

func (p Vec2) Floor() Vec2 {
	return Vec2{
		math.Floor(p.X),
		math.Floor(p.Y),
	}
}

func (p Vec2) Ceil() Vec2 {
	return Vec2{
		math.Ceil(p.X),
		math.Ceil(p.Y),
	}
}

func (p Vec2) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Vec2) LenSquared() float64 {
	return p.Dot(p)
}

func (p Vec2) Normalize() Vec2 {
	l := p.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{p.X / l, p.Y / l}
}

func (p Vec2) Scale2(q Vec2) Vec2 {
	return Vec2{p.X * q.X, p.Y * q.Y}
}

func (p Vec2) Scale(k float64) Vec2 {
	return Vec2{p.X * k, p.Y * k}
}

func (p Vec2) Shear(k float64) Vec2 {
	return Vec2{p.X + k*p.Y, p.Y + k*p.X}
}

func (p Vec2) Shearv(q Vec2) Vec2 {
	return Vec2{p.X + q.X*p.Y, p.Y + q.Y*p.X}
}

func (p Vec2) Lerp(t float64, q Vec2) Vec2 {
	return Vec2{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
	}
}

func (p Vec2) Lerp2(t, q Vec2) Vec2 {
	return Vec2{
		Lerp(t.X, p.X, q.X),
		Lerp(t.Y, p.Y, q.Y),
	}
}

func (p Vec2) Wrap(s, e float64) Vec2 {
	p.X = Wrap(p.X, s, e)
	p.Y = Wrap(p.Y, s, e)
	return p
}

func (p Vec2) Distance(q Vec2) float64 {
	return p.Sub(q).Len()
}

func (p Vec2) DistanceSquared(q Vec2) float64 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec2) Polar() Polar {
	return Polar{p.Len(), math.Atan2(p.Y, p.X)}
}

func (p Vec2) MinScalar(k float64) Vec2 {
	return Vec2{
		math.Min(p.X, k),
		math.Min(p.Y, k),
	}
}

func (p Vec2) MaxScalar(k float64) Vec2 {
	return Vec2{
		math.Max(p.X, k),
		math.Max(p.Y, k),
	}
}

func (p Vec2) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (p Vec2) Rotate(r float64) Vec2 {
	si, co := math.Sincos(r)
	return Vec2{
		p.X*co - p.Y*si,
		p.X*si + p.Y*co,
	}
}

func (p Vec2) Shrink(k float64) Vec2 {
	return Vec2{p.X / k, p.Y / k}
}

func (p Vec2) Shrink2(q Vec2) Vec2 {
	return Vec2{p.X / q.X, p.Y / q.Y}
}

func (p Vec2) YX() Vec2 {
	return Vec2{p.Y, p.X}
}

func (p Vec2) Clamp(s, e float64) Vec2 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	return p
}

func (p Vec2) Clamp2(s, e Vec2) Vec2 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	return p
}

func (p Vec2) Equals(q Vec2, eps float64) bool {
	return math.Abs(p.X-q.X) <= eps && math.Abs(p.Y-q.Y) <= eps
}

func (p Vec2) OnLine(a, b Vec2) bool {
	sx := math.Min(a.X, b.X)
	sy := math.Min(a.Y, b.Y)
	ex := math.Max(a.X, b.X)
	ey := math.Max(a.Y, b.Y)
	return sx <= p.X && p.X <= ex &&
		sy <= p.Y && p.Y <= ey
}

func (p Vec2) Angle(q Vec2) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
}

func (p Vec2) Point3() Vec3 {
	return Vec3{p.X, p.Y, 1}
}

func (p Vec2) Vec3() Vec3 {
	return Vec3{p.X, p.Y, 0}
}

func (p Vec2) String() string {
	return fmt.Sprintf(`Vec2(%0.3f, %0.3f)`, p.X, p.Y)
}

type Vec3 struct {
	X, Y, Z float64
}

func (p Vec3) Add(q Vec3) Vec3 {
	return Vec3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Vec3) Sub(q Vec3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Vec3) AddScale(q Vec3, t float64) Vec3 {
	return p.Add(q.Scale(t))
}

func (p Vec3) SubScale(q Vec3, t float64) Vec3 {
	return p.Sub(q.Scale(t))
}

func (p Vec3) Dot(q Vec3) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p Vec3) Cross(q Vec3) Vec3 {
	return Vec3{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

func (p Vec3) CrossNormalize(q Vec3) Vec3 {
	return p.Cross(q).Normalize()
}

func (p Vec3) Neg() Vec3 {
	return Vec3{-p.X, -p.Y, -p.Z}
}

func (p Vec3) Reflect(q Vec3) Vec3 {
	q = q.Scale(2 * p.Dot(q))
	return p.Sub(q)
}

func (p Vec3) Refract(q Vec3, eta float64) Vec3 {
	x := p.Dot(q)
	k := 1 - eta*eta*(1-x*x)
	if k < 0 {
		return Vec3{}
	}
	a := q.Scale(eta)
	b := p.Scale(eta*x + math.Sqrt(k))
	return a.Sub(b)
}

func (p Vec3) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Vec3) LenSquared() float64 {
	return p.Dot(p)
}

func (p Vec3) Scale3(q Vec3) Vec3 {
	return Vec3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

func (p Vec3) Scale(k float64) Vec3 {
	return Vec3{p.X * k, p.Y * k, p.Z * k}
}

func (p Vec3) Shrink(k float64) Vec3 {
	return Vec3{
		p.X / k,
		p.Y / k,
		p.Z / k,
	}
}

func (p Vec3) Shrink3(q Vec3) Vec3 {
	return Vec3{
		p.X / q.X,
		p.Y / q.Y,
		p.Z / q.Z,
	}
}

func (p Vec3) Normalize() Vec3 {
	l := p.Len()
	if l == 0 {
		return Vec3{}
	}
	return Vec3{p.X / l, p.Y / l, p.Z / l}
}

func (p Vec3) Distance(q Vec3) float64 {
	return p.Sub(q).Len()
}

func (p Vec3) DistanceSquared(q Vec3) float64 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec3) Lerp(t float64, q Vec3) Vec3 {
	return Vec3{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
		Lerp(t, p.Z, q.Z),
	}
}

func (p Vec3) Abs() Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
	}
}

func (p Vec3) MaxScalar(k float64) Vec3 {
	return Vec3{
		math.Max(p.X, k),
		math.Max(p.Y, k),
		math.Max(p.Z, k),
	}
}

func (p Vec3) MinScalar(k float64) Vec3 {
	return Vec3{
		math.Min(p.X, k),
		math.Min(p.Y, k),
		math.Min(p.Z, k),
	}
}

func (p Vec3) Max(q Vec3) Vec3 {
	return Vec3{
		math.Max(p.X, q.X),
		math.Max(p.Y, q.Y),
		math.Max(p.Z, q.Z),
	}
}

func (p Vec3) Min(q Vec3) Vec3 {
	return Vec3{
		math.Min(p.X, q.X),
		math.Min(p.Y, q.Y),
		math.Min(p.Z, q.Z),
	}
}

func (p Vec3) Point4() Vec4 {
	return Vec4{p.X, p.Y, p.Z, 1}
}

func (p Vec3) Vec4() Vec4 {
	return Vec4{p.X, p.Y, p.Z, 0}
}

func (p Vec3) MinComp() float64 {
	return math.Min(p.X, math.Min(p.Y, p.Z))
}

func (p Vec3) MaxComp() float64 {
	return math.Max(p.X, math.Max(p.Y, p.Z))
}

func (p Vec3) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	return color.RGBA{
		uint8(Clamp(p.X+0.5, 0, 255)),
		uint8(Clamp(p.Y+0.5, 0, 255)),
		uint8(Clamp(p.Z+0.5, 0, 255)),
		255,
	}

}

func (p Vec3) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec3) Spherical() Spherical {
	l := p.Len()
	return Spherical{
		R: l,
		T: math.Acos(p.Z / l),
		P: math.Atan2(p.Y, p.X),
	}
}

func (p Vec3) Equals(q Vec3, eps float64) bool {
	return math.Abs(p.X-q.X) <= eps && math.Abs(p.Y-q.Y) <= eps &&
		math.Abs(p.Z-q.Z) <= eps
}

func (p Vec3) XY() Vec2 { return Vec2{p.X, p.Y} }
func (p Vec3) XZ() Vec2 { return Vec2{p.X, p.Z} }
func (p Vec3) YX() Vec2 { return Vec2{p.Y, p.X} }
func (p Vec3) YZ() Vec2 { return Vec2{p.Y, p.Z} }
func (p Vec3) ZX() Vec2 { return Vec2{p.Z, p.X} }
func (p Vec3) ZY() Vec2 { return Vec2{p.Z, p.Y} }

func (p Vec3) Clamp(s, e float64) Vec3 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	p.Z = Clamp(p.Z, s, e)
	return p
}

func (p Vec3) Clamp3(s, e Vec3) Vec3 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	p.Z = Clamp(p.Z, s.Z, e.Z)
	return p
}

func (p Vec3) Angle(q Vec3) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
}

func (p Vec3) String() string {
	return fmt.Sprintf(`Vec3(%0.3f, %0.3f, %0.3f)`, p.X, p.Y, p.Z)
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (p Vec4) Add(q Vec4) Vec4 {
	return Vec4{p.X + q.X, p.Y + q.Y, p.Z + q.Z, p.W}
}

func (p Vec4) Sub(q Vec4) Vec4 {
	return Vec4{p.X - q.X, p.Y - q.Y, p.Z - q.Z, p.W}
}

func (p Vec4) AddScale(q Vec4, k float64) Vec4 {
	return p.Add(q.Scale(k))
}

func (p Vec4) SubScale(q Vec4, k float64) Vec4 {
	return p.Sub(q.Scale(k))
}

func (p Vec4) Scale(k float64) Vec4 {
	return Vec4{p.X * k, p.Y * k, p.Z * k, k * p.W}
}

func (p Vec4) Scale4(q Vec4) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W * q.W}
}

func (p Vec4) Shrink(k float64) Vec4 {
	return Vec4{p.X / k, p.Y / k, p.Z / k, p.W / k}
}

func (p Vec4) Shrink4(q Vec4) Vec4 {
	return Vec4{p.X / q.X, p.Y / q.Y, p.Z / q.Z, q.W / p.W}
}

func (p Vec4) Dot(q Vec4) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

func (p Vec4) Len() float64 {
	return math.Sqrt(p.Dot(p))
}

func (p Vec4) XYZ() Vec3 { return Vec3{p.X, p.Y, p.Z} }
func (p Vec4) XZY() Vec3 { return Vec3{p.X, p.Z, p.Y} }
func (p Vec4) YXZ() Vec3 { return Vec3{p.Y, p.X, p.Z} }
func (p Vec4) YZX() Vec3 { return Vec3{p.Y, p.Z, p.X} }

func (p Vec4) Normalize() Vec4 {
	l := p.Len()
	if l == 0 {
		return Vec4{0, 0, 0, p.W}
	}
	return Vec4{
		p.X / l,
		p.Y / l,
		p.Z / l,
		p.W / l,
	}
}

func (p Vec4) Floor() Vec4 {
	return Vec4{
		math.Floor(p.X),
		math.Floor(p.Y),
		math.Floor(p.Z),
		math.Floor(p.W),
	}
}

func (p Vec4) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	if 0 <= p.W && p.W <= 1 {
		p.W *= 255
	}
	c := color.RGBA{
		uint8(Clamp(p.X+0.5, 0, 255)),
		uint8(Clamp(p.Y+0.5, 0, 255)),
		uint8(Clamp(p.Z+0.5, 0, 255)),
		uint8(Clamp(p.W+0.5, 0, 255)),
	}
	return c
}

func (p Vec4) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec4) Angle(q Vec4) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
}

func (p Vec4) PerspectiveDivide() Vec4 {
	if p.W == 0 {
		return p
	}
	return Vec4{
		p.X / p.W,
		p.Y / p.W,
		p.Z / p.W,
		1,
	}
}

func (p Vec4) String() string {
	return fmt.Sprintf(`Vec4(%0.3f, %0.3f, %0.3f, %0.3f)`, p.X, p.Y, p.Z, p.W)
}

type Mat2 [2][2]float64

func (m *Mat2) Mul(a, b *Mat2) *Mat2 {
	var p Mat2
	for i := range a {
		for j := range a[i] {
			for k := range a[j] {
				p[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	*m = p
	return m
}

func (m *Mat2) Transform(p Vec2) Vec2 {
	return Vec2{
		m[0][0]*p.X + m[0][1]*p.Y,
		m[1][0]*p.X + m[1][1]*p.Y,
	}
}

func (m *Mat2) Transpose() *Mat2 {
	var p Mat2
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat2) Row(n int) Vec2 {
	return Vec2{m[n][0], m[n][1]}
}

func (m *Mat2) Col(n int) Vec2 {
	return Vec2{m[0][n], m[1][n]}
}

func (m *Mat2) Trace() float64 {
	return m[0][0] + m[1][1]
}

func (m *Mat2) Det() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

func (m *Mat2) Inverse() *Mat2 {
	invdet := 1 / m.Det()
	fmt.Println(m.Det())
	var minv Mat2
	minv[0][0] = m[1][1] * invdet
	minv[0][1] = -m[0][1] * invdet
	minv[1][0] = -m[1][0] * invdet
	minv[1][1] = m[0][0] * invdet

	*m = minv
	return m
}

func (m Mat2) String() string {
	return fmt.Sprintf(`
Mat2[% 0.3f, % 0.3f,
     % 0.3f, % 0.3f]`,
		m[0][0], m[0][1],
		m[1][0], m[1][1])
}

type Mat3 [3][3]float64

func (m *Mat3) Identity() *Mat3 {
	*m = Mat3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	return m
}

func (m *Mat3) Mul(a, b *Mat3) *Mat3 {
	var p Mat3
	for i := range a {
		for j := range a[i] {
			for k := range a[j] {
				p[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	*m = p
	return m
}

func (m *Mat3) Transform(p Vec3) Vec3 {
	return Vec3{
		m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z,
		m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z,
		m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z,
	}
}

func (m *Mat3) Transpose() *Mat3 {
	var p Mat3
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat3) Det() float64 {
	A1 := Vec3{m[0][0], m[0][1], m[0][2]}
	A2 := Vec3{m[1][1], m[1][2], m[1][0]}
	A3 := Vec3{m[2][2], m[2][0], m[2][1]}
	A4 := Vec3{m[1][2], m[1][0], m[1][1]}
	A5 := Vec3{m[2][1], m[2][2], m[2][0]}

	X := A2.Scale3(A3)
	Y := A4.Scale3(A5)
	A6 := X.Sub(Y)
	return A1.Dot(A6)
}

func (m *Mat3) Trace() float64 {
	return m[0][0] + m[1][1] + m[2][2]
}

func (m *Mat3) Mat4() Mat4 {
	return Mat4{
		{m[0][0], m[1][0], m[2][0], 0},
		{m[0][1], m[1][1], m[2][1], 0},
		{m[0][2], m[1][2], m[2][2], 0},
		{0, 0, 0, 1},
	}
}

func (m *Mat3) FromBasis(X, Y, Z Vec3) *Mat3 {
	*m = Mat3{
		{X.X, Y.X, Z.X},
		{X.Y, Y.Y, Z.Y},
		{X.Z, Y.Z, Z.Z},
	}
	return m
}

func (m *Mat3) Basis() (X, Y, Z, W Vec3) {
	X = Vec3{m[0][0], m[1][0], m[2][0]}
	Y = Vec3{m[0][1], m[1][1], m[2][1]}
	Z = Vec3{m[0][2], m[1][2], m[2][2]}
	return
}

func (m *Mat3) SetCol(n int, p Vec3) {
	m[0][n] = p.X
	m[1][n] = p.Y
	m[2][n] = p.Z
}

func (m *Mat3) SetRow(n int, p Vec3) {
	m[n][0] = p.X
	m[n][1] = p.Y
	m[n][2] = p.Z
}

func (m *Mat3) Row(n int) Vec3 {
	return Vec3{m[n][0], m[n][1], m[n][2]}
}

func (m *Mat3) Col(n int) Vec3 {
	return Vec3{m[0][n], m[1][n], m[2][n]}
}

func (m *Mat3) Inverse() *Mat3 {
	det := m[0][0]*(m[1][1]*m[2][2]-m[2][1]*m[1][2]) -
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
	invdet := 1 / det

	var minv Mat3
	minv[0][0] = (m[1][1]*m[2][2] - m[2][1]*m[1][2]) * invdet
	minv[0][1] = (m[0][2]*m[2][1] - m[0][1]*m[2][2]) * invdet
	minv[0][2] = (m[0][1]*m[1][2] - m[0][2]*m[1][1]) * invdet
	minv[1][0] = (m[1][2]*m[2][0] - m[1][0]*m[2][2]) * invdet
	minv[1][1] = (m[0][0]*m[2][2] - m[0][2]*m[2][0]) * invdet
	minv[1][2] = (m[1][0]*m[0][2] - m[0][0]*m[1][2]) * invdet
	minv[2][0] = (m[1][0]*m[2][1] - m[2][0]*m[1][1]) * invdet
	minv[2][1] = (m[2][0]*m[0][1] - m[0][0]*m[2][1]) * invdet
	minv[2][2] = (m[0][0]*m[1][1] - m[1][0]*m[0][1]) * invdet
	*m = minv
	return m
}

func (m Mat3) String() string {
	return fmt.Sprintf(`
Mat3[% 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}

type Mat4 [4][4]float64

func (m *Mat4) Identity() *Mat4 {
	*m = Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Add(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func (m *Mat4) Sub(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func (m *Mat4) Mul(a, b *Mat4) *Mat4 {
	var p Mat4
	for i := range a {
		for j := range a[i] {
			for k := range a[j] {
				p[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	*m = p
	return m
}

func (m *Mat4) Trace() float64 {
	return m[0][0] + m[1][1] + m[2][2] + m[3][3]
}

func (m *Mat4) Translate3(p Vec3) *Mat4 {
	m.Translate(p.X, p.Y, p.Z)
	return m
}

func (m *Mat4) Translate(tx, ty, tz float64) *Mat4 {
	*m = Mat4{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Scale3(p Vec3) *Mat4 {
	m.Scale(p.X, p.Y, p.Z)
	return m
}

func (m *Mat4) Scale(sx, sy, sz float64) *Mat4 {
	*m = Mat4{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) LookAt(eye, center, up Vec3) *Mat4 {
	f := center
	f = f.Sub(eye)
	f = f.Normalize()

	s := f.Cross(up)
	s = s.Normalize()
	u := s.Cross(f)

	*m = Mat4{
		{s.X, s.Y, s.Z, 0},
		{u.X, u.Y, u.Z, 0},
		{-f.X, -f.Y, -f.Z, 0},
		{0, 0, 0, 1},
	}

	var t Mat4
	t.Translate(-eye.X, -eye.Y, -eye.Z)
	m.Mul(&t, m)
	return m
}

func (m *Mat4) RotateX(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{1, 0, 0, 0},
		{0, co, -si, 0},
		{0, si, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotateY(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{co, 0, si, 0},
		{0, 1, 0, 0},
		{-si, 0, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotateZ(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{co, -si, 0, 0},
		{si, co, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Frustum(l, r, b, t, n, f float64) *Mat4 {
	A := (r + l) / (r - l)
	B := (t + b) / (t - b)
	C := -(f + n) / (f - n)
	D := -2 * f * n / (f - n)
	*m = Mat4{
		{2 * n / (r - l), 0, A, 0},
		{0, 2 * n / (t - b), B, 0},
		{0, 0, C, D},
		{0, 0, -1, 0},
	}
	return m
}

func (m *Mat4) InfinitePerspective(fovy, aspect, near float64) *Mat4 {
	const zp = 0
	f := 1 / math.Tan(fovy/2)
	*m = Mat4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, -(1 - zp), -near * (1 - zp)},
		{0, 0, -1, 0},
	}
	return m
}

func (m *Mat4) Perspective(fovy, aspect, near, far float64) *Mat4 {
	ymax := near * math.Tan(fovy/2)
	xmax := ymax * aspect
	m.Frustum(-xmax, xmax, -ymax, ymax, near, far)
	return m
}

func (m *Mat4) Ortho(l, r, b, t, n, f float64) *Mat4 {
	sx := 2 / (r - l)
	sy := 2 / (t - b)
	sz := -2 / (f - n)

	tx := -(r + l) / (r - l)
	ty := -(t + b) / (t - b)
	tz := -(f + n) / (f - n)

	*m = Mat4{
		{sx, 0, 0, tx},
		{0, sy, 0, ty},
		{0, 0, sz, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Viewport(x, y, w, h float64) *Mat4 {
	l := x
	b := y
	r := x + w
	t := y + h
	*m = Mat4{
		{(r - l) / 2, 0, 0, (r + l) / 2},
		{0, (t - b) / 2, 0, (t + b) / 2},
		{0, 0, 0.5, 0.5},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Inverse() *Mat4 {
	a := Vec3{m[0][0], m[1][0], m[2][0]}
	b := Vec3{m[0][1], m[1][1], m[2][1]}
	c := Vec3{m[0][2], m[1][2], m[2][2]}
	d := Vec3{m[0][3], m[1][3], m[2][3]}

	s := a.Cross(b)
	t := c.Cross(d)

	invDet := 1 / s.Dot(c)

	s = s.Scale(invDet)
	t = t.Scale(invDet)
	v := c.Scale(invDet)

	r0 := b.Cross(v)
	r1 := v.Cross(a)

	*m = Mat4{
		{r0.X, r0.Y, r0.Z, -b.Dot(t)},
		{r1.X, r1.Y, r1.Z, a.Dot(t)},
		{s.X, s.Y, s.Z, -d.Dot(s)},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Transpose() *Mat4 {
	var p Mat4
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat4) Transform(v Vec4) Vec4 {
	return Vec4{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}

func (m *Mat4) Transform3(v Vec3) Vec3 {
	s := m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]
	switch s {
	case 0:
		return Vec3{}
	default:
		invs := 1 / s
		p := m.Transform(Vec4{v.X, v.Y, v.Z, 1})
		return Vec3{p.X * invs, p.Y * invs, p.Z * invs}
	}

}

func (m *Mat4) FromBasis3(X, Y, Z, W Vec3) *Mat4 {
	*m = Mat4{
		{X.X, Y.X, Z.X, 0},
		{X.Y, Y.Y, Z.Y, 0},
		{X.Z, Y.Z, Z.Z, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) FromBasis(X, Y, Z, W Vec4) *Mat4 {
	*m = Mat4{
		{X.X, Y.X, Z.X, W.X},
		{X.Y, Y.Y, Z.Y, W.Y},
		{X.Z, Y.Z, Z.Z, W.Z},
		{X.W, Y.W, Z.W, W.W},
	}
	return m
}

func (m *Mat4) Basis3() (X, Y, Z, W Vec3) {
	X = Vec3{m[0][0], m[1][0], m[2][0]}
	Y = Vec3{m[0][1], m[1][1], m[2][1]}
	Z = Vec3{m[0][2], m[1][2], m[2][2]}
	W = Vec3{m[0][3], m[1][3], m[2][3]}
	return
}

func (m *Mat4) Basis() (X, Y, Z, W Vec4) {
	X = Vec4{m[0][0], m[1][0], m[2][0], m[3][0]}
	Y = Vec4{m[0][1], m[1][1], m[2][1], m[3][1]}
	Z = Vec4{m[0][2], m[1][2], m[2][2], m[3][2]}
	W = Vec4{m[0][3], m[1][3], m[2][3], m[3][3]}
	return
}

func (m *Mat4) SetCol(n int, p Vec4) {
	m[0][n] = p.X
	m[1][n] = p.Y
	m[2][n] = p.Z
	m[3][n] = p.W
}

func (m *Mat4) SetRow(n int, p Vec4) {
	m[n][0] = p.X
	m[n][1] = p.Y
	m[n][2] = p.Z
	m[n][3] = p.W
}

func (m *Mat4) Row(n int) Vec4 {
	return Vec4{m[n][0], m[n][1], m[n][2], m[n][3]}
}

func (m *Mat4) Col(n int) Vec4 {
	return Vec4{m[0][n], m[1][n], m[2][n], m[3][n]}
}

func (m Mat4) String() string {
	return fmt.Sprintf(`
Mat4[% 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3])
}

type Polar struct {
	R, P float64
}

func (p *Polar) Cartesian() Vec2 {
	s, c := math.Sincos(p.P)
	return Vec2{p.R * c, p.R * s}
}

type Quat struct {
	X, Y, Z, W float64
}

func (q Quat) Add(r Quat) Quat {
	return Quat{
		q.X + r.X,
		q.Y + r.Y,
		q.Z + r.Z,
		q.W + r.W,
	}
}

func (q Quat) Sub(r Quat) Quat {
	return Quat{
		q.X - r.X,
		q.Y - r.Y,
		q.Z - r.Z,
		q.W - r.W,
	}
}

func (q Quat) Mul(r Quat) Quat {
	return Quat{
		q.X*r.W + q.Y*r.Z - q.Z*r.Y + q.W*r.X,
		q.Y*r.W + q.Z*r.X + q.W*r.Y - q.X*r.Z,
		q.Z*r.W + q.W*r.Z + q.X*r.Y - q.Y*r.X,
		q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z,
	}
}

func (q Quat) Neg() Quat {
	return Quat{-q.X, -q.Y, -q.Z, -q.W}
}

func (q Quat) Dot(p Quat) float64 {
	return q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
}

func (q Quat) Scale(k float64) Quat {
	return Quat{
		q.X * k,
		q.Y * k,
		q.Z * k,
		q.W * k,
	}
}

func (q Quat) Len() float64 {
	return math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
}

func (q Quat) Normalize() Quat {
	l := q.Len()
	if l == 0 {
		return Quat{}
	}
	return Quat{
		q.X / l,
		q.Y / l,
		q.Z / l,
		q.W / l,
	}
}

func (q Quat) Conj() Quat {
	return Quat{-q.X, -q.Y, -q.Z, q.W}
}

func (q Quat) FromAxisAngle(v Vec3, r float64) Quat {
	r *= 0.5
	vn := v.Normalize()
	si, co := math.Sincos(r)
	return Quat{
		vn.X * si,
		vn.Y * si,
		vn.Z * si,
		co,
	}
}

func (q Quat) FromEuler(pitch, yaw, roll float64) Quat {
	p := pitch / 2
	y := yaw / 2
	r := roll / 2

	sinp := math.Sin(p)
	siny := math.Sin(y)
	sinr := math.Sin(r)
	cosp := math.Cos(p)
	cosy := math.Cos(y)
	cosr := math.Cos(r)

	return Quat{
		sinr*cosp*cosy - cosr*sinp*siny,
		cosr*sinp*cosy + sinr*cosp*siny,
		cosr*cosp*siny - sinr*sinp*cosy,
		cosr*cosp*cosy + sinr*sinp*siny,
	}.Normalize()
}

func (q Quat) Mat3() Mat3 {
	x, y, z, w := q.X, q.Y, q.Z, q.W
	x2 := x * x
	y2 := y * y
	z2 := z * z
	xy := x * y
	xz := x * z
	yz := y * z
	wx := w * x
	wy := w * y
	wz := w * z

	return Mat3{
		{1.0 - 2.0*(y2+z2), 2.0 * (xy - wz), 2.0 * (xz + wy)},
		{2.0 * (xy + wz), 1.0 - 2.0*(x2+z2), 2.0 * (yz - wx)},
		{2.0 * (xz - wy), 2.0 * (yz + wx), 1.0 - 2.0*(x2+y2)},
	}
}

func (q Quat) Mat4() Mat4 {
	m := q.Mat3()
	return m.Mat4()
}

func (q Quat) Transform3(v Vec3) Vec3 {
	m := q.Mat3()
	return m.Transform(v)
}

func (q Quat) Transform(v Vec4) Vec4 {
	m := q.Mat4()
	return m.Transform(v)
}

func (q Quat) AxisAngle() (v Vec3, r float64) {
	s := math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	v = Vec3{q.X / s, q.Y / s, q.Z / s}
	r = math.Acos(q.W) * 2
	return
}

func (q Quat) Lerp(t float64, p Quat) Quat {
	return q.Add(p.Sub(q).Scale(t))
}

func (q Quat) Inverse() Quat {
	l2 := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	return Quat{
		-q.X / l2,
		-q.Y / l2,
		-q.Z / l2,
		q.W / l2,
	}
}

func (q Quat) Powu(p float64) Quat {
	t := math.Acos(q.W) * p
	u := Vec3{q.X, q.Y, q.Z}
	u = u.Normalize()
	u = u.Scale(math.Sin(t))
	w := math.Cos(t)
	return Quat{
		u.X, u.Y, u.Z, w,
	}
}

func (q Quat) Slerp(t float64, p Quat) Quat {
	v0 := q.Normalize()
	v1 := p.Normalize()

	const threshold = 0.9995
	dot := v0.Dot(v1)
	if dot > threshold {
		return v0.Lerp(t, v1).Normalize()
	}

	if dot < 0 {
		v1 = v1.Neg()
		dot = -dot
	}

	dot = Clamp(dot, -1, 1)
	theta0 := math.Acos(dot)
	theta := theta0 * t

	v2 := v1.Sub(v0.Scale(dot))
	v2 = v2.Normalize()

	v3 := v0.Scale(math.Cos(theta))
	v4 := v2.Scale(math.Sin(theta))
	return v3.Add(v4)
}

func (q Quat) String() string {
	return fmt.Sprintf(`Quat(%0.3f, %0.3f, %0.3f, %0.3f)`, q.X, q.Y, q.Z, q.W)
}

type Spherical struct {
	R, T, P float64
}

func (s Spherical) Euclidean() Vec3 {
	sint := math.Sin(s.T)
	cost := math.Sin(s.T)
	sinp := math.Sin(s.P)
	cosp := math.Sin(s.P)
	r := s.R

	return Vec3{
		r * sint * cosp,
		r * sint * sinp,
		r * cost,
	}
}

func Lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func Unlerp(t, a, b float64) float64 {
	return (t - a) / (b - a)
}

func LinearRemap(x, a, b, c, d float64) float64 {
	return Lerp(Unlerp(x, a, b), c, d)
}

func Smoothstep(a, b, x float64) float64 {
	t := Clamp((x-a)/(b-a), 0, 1)
	return t * t * (3 - 2*t)
}

func CubicBezier1D(t, p0, p1, p2, p3 float64) float64 {
	it := 1 - t
	return it*it*it*p0 + 3*it*it*t*p1 + 3*it*t*t*p2 + t*t*t*p3
}

func Clamp(x, s, e float64) float64 {
	if x < s {
		x = s
	}
	if x > e {
		x = e
	}
	return x
}

func Saturate(x float64) float64 {
	return math.Max(0, math.Min(1, x))
}

func SignStrict(x float64) float64 {
	if x >= 0 {
		return 1
	}
	return -1
}

func Sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	if x == 0 {
		return 0
	}
	return 1
}

func Deg2Rad(d float64) float64 {
	return d * math.Pi / 180
}

func Rad2Deg(r float64) float64 {
	return r * 180 / math.Pi
}

type Circle struct {
	X, Y, R float64
}

func (c Circle) InPoint(x, y float64) bool {
	return (x-c.X)*(x-c.X)+(y-c.Y)*(y-c.Y) <= c.R
}

func (c Circle) InRect(r Rectangle) bool {
	dx := c.X - math.Max(r.Min.X, math.Min(c.X, r.Max.X))
	dy := c.Y - math.Max(r.Min.Y, math.Min(c.Y, r.Max.Y))
	return dx*dx+dy*dy <= c.R
}

type Rectangle struct {
	Min, Max Vec2
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}

func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}

	if r.Empty() {
		return Rectangle{}
	}
	return r
}

func (r Rectangle) Union(s Rectangle) Rectangle {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

func (r Rectangle) Int() image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r Rectangle) Scale(s Vec2) Rectangle {
	r.Min.X *= s.X
	r.Max.X *= s.X
	r.Min.Y *= s.Y
	r.Max.Y *= s.Y
	return r
}

func (r Rectangle) Add(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X + p.X, r.Min.Y + p.Y},
		Vec2{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

func (r Rectangle) Sub(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X - p.X, r.Min.Y - p.Y},
		Vec2{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

func (r Rectangle) Size() Vec2 {
	return Vec2{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

func (r Rectangle) Floor() Rectangle {
	return Rectangle{
		r.Min.Floor(),
		r.Max.Floor(),
	}
}

func (r Rectangle) Inset(n float64) Rectangle {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Center() Vec2 {
	return Vec2{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rectangle) Diagonal() float64 {
	x := r.Max.X - r.Min.X
	y := r.Max.Y - r.Min.Y
	return math.Sqrt(x*x + y*y)
}

func (r Rectangle) PosSize() (x, y, w, h float64) {
	return r.Min.X, r.Min.Y, r.Dx(), r.Dy()
}

func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}

	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

func (r Rectangle) TL() Vec2 { return r.Min }
func (r Rectangle) TR() Vec2 { return Vec2{r.Max.X, r.Min.Y} }
func (r Rectangle) BL() Vec2 { return Vec2{r.Min.X, r.Max.Y} }
func (r Rectangle) BR() Vec2 { return r.Max }

func (r Rectangle) Inverted() bool {
	return r.Min.X > r.Max.X || r.Min.Y > r.Max.Y
}

func (r Rectangle) Expand(x, y float64) Rectangle {
	r.Min.X -= x
	r.Min.Y -= y
	r.Max.X += x
	r.Max.Y += y
	return r
}

func (r Rectangle) Expand2(v Vec2) Rectangle {
	return r.Expand(v.X, v.Y)
}

func RoundPrec(v float64, prec int) float64 {
	if prec < 0 {
		return v
	}

	tab := [...]float64{
		1, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8, 1e-9, 1e-10,
	}
	step := 0.0
	if prec < len(tab) {
		step = tab[prec]
	} else {
		step = math.Pow(10, float64(-prec))
	}

	neg := v < 0
	v = math.Abs(v)
	rem := math.Mod(v, step)
	if rem <= step*0.5 {
		v -= rem
	} else {
		v += step - rem
	}

	if neg {
		v = -v
	}

	return v
}

func Sinc(x float64) float64 {
	x *= math.Pi
	if x < 0.01 && x > -0.01 {
		return 1 + x*x*((-1.0/6)+x*x*1.0/120)
	}
	return math.Sin(x) / x
}

func LinearController(curpos *float64, targetpos, acc, deacc, dt float64) {
	sign := 1.0
	p := 0.0
	cp := *curpos
	if cp == targetpos {
		return
	}
	if targetpos < cp {
		targetpos = -targetpos
		cp = -cp
		sign = -1
	}

	// first decelerate
	if cp < 0 {
		p = cp + deacc*dt
		if p > 0 {
			p = 0
			dt = dt - p/deacc
			if dt < 0 {
				dt = 0
			}
		} else {
			dt = 0
		}
		cp = p
	}

	// now accelerate
	p = cp + acc*dt
	if p > targetpos {
		p = targetpos
	}
	*curpos = p * sign
}

func Multiple(a, m float64) float64 {
	return math.Ceil(a/m) * m
}

func Wrap(x, s, e float64) float64 {
	if x < s {
		x += e
	}
	if x >= e {
		x -= e
	}
	return x
}

func TriangleBarycentric(p, a, b, c Vec2) Vec3 {
	x := Vec3{
		c.X - a.X,
		b.X - a.X,
		a.X - p.X,
	}
	y := Vec3{
		c.Y - a.Y,
		b.Y - a.Y,
		a.Y - p.Y,
	}
	u := x.Cross(y)
	if math.Abs(u.Z) > 1e-2 {
		return Vec3{
			1 - (u.X+u.Y)/u.Z,
			u.Y / u.Z,
			u.X / u.Z,
		}
	}
	return Vec3{-1, -1, -1}
}

func Clamp8(x, a, b float64) uint8 {
	x = math.Round(x)
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return uint8(x)
}

func ditfft2c(x, y []complex128, n, s int) {
	if n == 1 {
		y[0] = x[0]
		return
	}
	ditfft2c(x, y, n/2, 2*s)
	ditfft2c(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

func ditfft2r(x []float64, y []complex128, n, s int) {
	if n == 1 {
		y[0] = complex(x[0], 0)
		return
	}
	ditfft2r(x, y, n/2, 2*s)
	ditfft2r(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

func FFT1DC(dst, src []complex128) {
	ditfft2c(src, dst, len(dst), 1)
}

func IFFT1DC(dst, src []complex128) {
	for i := range src {
		src[i] = cmplx.Conj(src[i])
	}
	FFT1DC(dst, src)
	for i := range src {
		src[i] = cmplx.Conj(src[i])
	}
	for i := range dst {
		dst[i] = cmplx.Conj(dst[i]) / complex(float64(len(dst)), 0)
	}
}

func FFT1DR(dst []complex128, src []float64) {
	ditfft2r(src, dst, len(dst), 1)
}

func IFFT1DR(dst []float64, src []complex128) {
	for i := range src {
		src[i] = cmplx.Conj(src[i])
	}
	tmp := make([]complex128, len(dst))
	FFT1DC(tmp, src)
	for i := range src {
		src[i] = cmplx.Conj(src[i])
	}
	for i := range dst {
		dst[i] = real(tmp[i]) / float64(len(dst))
	}
}

func Simpson1D(f func(x float64) float64, start, end float64, n int) float64 {
	r := 0.0
	s := (end - start) / float64(n)
	i := 0

	r += f(start)
	for j := 1; j < n; j++ {
		r += (4 - float64(i<<1)) * f(start+float64(j)*s)
		i = (i + 1) & 1
	}
	r += f(end)
	r *= s / 3
	return r
}

func simpsonweight(i, n int) float64 {
	if i == 0 || i == n {
		return 1
	}
	if i%2 != 0 {
		return 4
	}
	return 2
}

func Simpson2D(f func(x, y float64) float64, x0, x1, y0, y1 float64, m, n int) float64 {
	if n%2 != 0 || m%2 != 0 {
		panic("integration range must be even")
	}

	dx := (x1 - x0) / float64(m)
	dy := (y1 - y0) / float64(n)
	r := 0.0
	for i := 0; i <= n; i++ {
		y := y0 + float64(i)*dy
		wy := simpsonweight(i, n)
		for j := 0; j <= m; j++ {
			x := x0 + float64(j)*dx
			wx := simpsonweight(j, m)
			r += f(x, y) * wx * wy
		}
	}
	r *= dx * dy / (9 * float64(m) * float64(n))
	return r
}

func Convolve1D(dst, src, coeffs []float64, shape int) []float64 {
	var m, n int
	switch shape {
	case 'f':
		m = len(src) + len(coeffs) - 1
	case 's':
		m = len(src)
		n = len(coeffs) - 2
	case 'v':
		m = len(src) - len(coeffs) + 1
		n = len(coeffs) - 1
		if m < 0 {
			m = 0
		}
	default:
		panic("unsupported convolution shape")
	}

	for k := 0; k < m; k++ {
		dst[k] = 0
		for j := range src {
			l := k + n - j
			if l < 0 || l >= len(coeffs) {
				continue
			}
			dst[k] += src[j] * coeffs[l]
		}
	}

	return dst[:m]
}

func Sample1D(f func(float64) float64, x0, x1 float64, n int) (p []float64, s float64) {
	p = make([]float64, n)
	s = (x1 - x0) / float64(n-1)
	for i := 0; i < n; i++ {
		p[i] = f(x0 + float64(i)*s)
	}
	return
}

func Sample2D(f func(x, y float64) float64, x0, x1, y0, y1 float64, nx, ny int) (p []float64, sx, sy float64) {
	p = make([]float64, nx*ny)
	sx = (x1 - x0) / float64(nx-1)
	sy = (y1 - y0) / float64(ny-1)
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			p[y*nx+x] = f(x0+sx*float64(x), y0+sy*float64(y))
		}
	}
	return
}

func FloatToComplex(v []float64) []complex128 {
	p := make([]complex128, len(v))
	for i := range v {
		p[i] = complex(v[i], 0)
	}
	return p
}

func ComplexToFloat(v []complex128) []float64 {
	p := make([]float64, len(v))
	for i := range v {
		p[i] = cmplx.Abs(v[i])
	}
	return p
}

func PearsonCorrelation1D(x, y []float64) float64 {
	return Cov1D(x, y, 1) / (Stddev(x, 1) * Stddev(y, 1))
}

func Cov1D(x, y []float64, ddof float64) float64 {
	mx := Mean(x)
	my := Mean(y)
	s := 0.0
	for i := range x {
		s += (x[i] - mx) * (y[i] * my)
	}
	return s / (float64(len(x)) - ddof)
}

func Mean(x []float64) float64 {
	s := 0.0
	for i := range x {
		s += x[i]
	}
	return s / float64(len(x))
}

func Stddev(x []float64, ddof float64) float64 {
	if len(x) <= 1 {
		return 0
	}
	xm := Mean(x)
	s := 0.0
	for i := range x {
		s += (x[i] - xm) * (x[i] - xm)
	}
	return math.Sqrt(s / (float64(len(x)) - ddof))
}

func Median(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	sort.Float64s(x)
	return x[len(x)/2]
}

func Mins(x ...float64) float64 {
	if len(x) == 0 {
		return 0
	}

	n := x[0]
	for i := range x[1:] {
		n = math.Min(n, x[i])
	}
	return n
}

func Maxs(x ...float64) float64 {
	if len(x) == 0 {
		return 0
	}

	n := x[0]
	for i := range x[1:] {
		n = math.Max(n, x[i])
	}
	return n
}
