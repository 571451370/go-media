package imgui

import (
	"math"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) CalcTextSize(text string, hideTextAfterDoubleHash bool, wrapWidth float64) f64.Vec2 {
	if hideTextAfterDoubleHash {
		n := strings.Index(text, "##")
		if n >= 0 {
			text = text[:n]
		}
	}

	font := c.Font
	fontSize := c.FontSize
	if text == "" {
		return f64.Vec2{0, fontSize}
	}
	textSize, _ := font.CalcTextSizeA(fontSize, math.MaxFloat64, wrapWidth, text)

	// Cancel out character spacing for the last character of a line (it is baked into glyph->AdvanceX field)
	fontScale := fontSize / font.FontSize
	characterSpacingX := fontScale
	if textSize.X > 0 {
		textSize.X -= characterSpacingX
	}
	textSize.X = float64(int(textSize.X + 0.95))
	return textSize
}
