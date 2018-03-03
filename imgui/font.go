package imgui

import "github.com/qeedquan/go-media/math/f64"

type Font struct {
	FontSize        float64
	Scale           float64
	DisplayOffset   f64.Vec2
	Glyphs          []FontGlyph
	IndexAdvanceX   []float64
	IndexLookup     []int
	FallbackGlyph   *FontGlyph
	FallbackChar    rune
	Ascent, Descent float64
}

type FontGlyph struct {
}

type FontAtlas struct {
	Flags FontAtlasFlags // Build flags (see FontAtlasFlags)
	Fonts []*Font        // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.
}

type FontAtlasFlags uint

const (
	AtlasFlagsNoPowerOfTwoHeight FontAtlasFlags = 1 << iota // Don't round the height to next power of two
	AtlasFlagsNoMouseCursors                                // Don't build software mouse cursors into the atlas
)

func (c *Context) GetDefaultFont() *Font {
	io := c.GetIO()
	if io.FontDefault != nil {
		return io.FontDefault
	}
	return io.Fonts.Fonts[0]
}

func (c *Context) SetCurrentFont(font *Font) {
	io := c.GetIO()
	window := c.GetCurrentWindow()
	c.Font = font
	c.FontBaseSize = io.FontGlobalScale * c.Font.FontSize * c.Font.Scale
	c.FontSize = 0
	if window != nil {
		c.FontSize = window.CalcFontSize()
	}
}