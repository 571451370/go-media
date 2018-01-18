package resampler

import (
	"image"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

var (
	srgb   [256]float64
	linear [4096]int
)

func init() {
	const gamma = 1.75
	for i := range srgb {
		srgb[i] = math.Pow(float64(i)/255, gamma)
	}

	for i := range linear {
		k := 255 * (math.Pow(float64(i)/float64(len(linear)), 1/gamma) + .5)
		k = f64.Clamp(k, 0, 255)
		linear[i] = int(k)
	}
}

func ResizeImage(m image.Image, dr image.Rectangle) *image.RGBA {
	var (
		resamplers [4]*Resampler
		samples    [4][]float64
	)
	sr := m.Bounds()
	sn := image.Pt(sr.Dx(), sr.Dy())
	dn := image.Pt(dr.Dx(), dr.Dy())
	for i := range resamplers {
		resamplers[i] = New(sn, dn, nil)
		samples[i] = make([]float64, sn.X)
	}

	y := 0
	var n int
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			c := color.RGBAModel.Convert(m.At(x, y)).(color.RGBA)
			samples[0][n] = srgb[c.R]
			samples[1][n] = srgb[c.G]
			samples[2][n] = srgb[c.B]
			samples[3][n] = float64(c.A) / 255
		}

		for i, rp := range resamplers {
			rp.PutLine(samples[i])
		}

		for {
			i := 0
			for i = range resamplers {
				out := resamplers[i].GetLine()
				if out == nil {
					break
				}
			}

			if i < len(resamplers) {
				break
			}
			y++
		}
	}

	return nil
}
