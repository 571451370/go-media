package chroma

import (
	"image/color"
	"math"
)

func DistanceLPRGB(a, b color.Color, p float64) float64 {
	x := color.RGBAModel.Convert(a).(color.RGBA)
	y := color.RGBAModel.Convert(b).(color.RGBA)

	d1 := math.Pow(float64(x.R-y.R), p)
	d2 := math.Pow(float64(x.G-y.G), p)
	d3 := math.Pow(float64(x.B-y.B), p)
	return math.Pow(d1+d2+d3, 1/p)
}

func DistanceL2RGB(a, b color.Color) float64 {
	return DistanceLPRGB(a, b, 2)
}

// https://www.compuphase.com/cmetric.htm
func DistanceWL2RGB(a, b color.Color) float64 {
	x := color.RGBAModel.Convert(a).(color.RGBA)
	y := color.RGBAModel.Convert(b).(color.RGBA)

	rm := float64(x.R+y.R) / 2
	cr := float64(x.R - y.R)
	cg := float64(x.G - y.G)
	cb := float64(x.B - y.B)
	return math.Sqrt(((512+rm)*cr*cr)/256 + 4*cg*cg + ((767-rm)*cb*cb)/256)
}

func DistanceBW(a, b color.Color) float64 {
	x := color.GrayModel.Convert(a).(color.Gray)
	y := color.GrayModel.Convert(b).(color.Gray)
	if x.Y != y.Y {
		return 1
	}
	return 0
}
