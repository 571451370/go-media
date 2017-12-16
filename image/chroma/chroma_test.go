package chroma

import (
	"fmt"
	"image/color"
	"testing"
)

func TestChroma(t *testing.T) {
	c := color.RGBA{255, 235, 241, 255}
	h := HSVModel.Convert(c)
	fmt.Println(c, h)
}
