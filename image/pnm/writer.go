package pnm

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
)

type Options struct {
	Format int
}

func Encode(w io.Writer, m image.Image, o *Options) error {
	if o == nil {
		o = &Options{Format: 3}
	}

	b := bufio.NewWriter(w)
	r := m.Bounds()

	fmt.Fprintf(w, "P%d", o.Format)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			p := m.At(x, y)
			switch o.Format {
			case 1, 2:
				c := (color.GrayModel.Convert(p)).(color.Gray)
				if o.Format == 1 && c.Y != 0 {
					c.Y = 1
				}
				fmt.Fprintf(w, "%d ", c.Y)
			case 3:
				c := color.RGBAModel.Convert(p).(color.RGBA)
				fmt.Fprintf(w, "%d %d %d", c.R, c.G, c.B)
			default:
				return fmt.Errorf("unknown pnm format")
			}

			if x+1 < r.Max.X {
				fmt.Fprintf(w, " ")
			}
		}

		switch o.Format {
		case 1, 2, 3:
			fmt.Fprintf(w, "\n")
		}
	}

	return b.Flush()
}
