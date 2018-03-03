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
	Flags           FontAtlasFlags // Build flags (see FontAtlasFlags)
	TexId           TextureID      // User data to refer to the texture once it has been uploaded to user's graphic systems. It is passed back to you during rendering via the ImDrawCmd structure.
	TexDesiredWidth int            // Texture width desired by user before Build(). Must be a power-of-two. If have many glyphs your graphics API have texture size restrictions you may want to increase texture width to decrease height.
	TexGlyphPadding int            // Padding between glyphs within texture in pixels. Defaults to 1.

	Fonts []*Font // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.

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