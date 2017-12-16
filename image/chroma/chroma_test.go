package chroma

import (
	"fmt"
	"image/color"
	"testing"
)

func TestChroma(t *testing.T) {
	hsv := HSVModel.Convert(color.RGBA{34, 19, 255, 255}).(HSV)
	rgba := color.RGBAModel.Convert(hsv).(color.RGBA)
	fmt.Println(hsv, rgba)
}
