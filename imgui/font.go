package imgui

import (
	"math"
	"unicode/utf8"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type Font struct {
	// Members: Hot ~62/78 bytes
	FontSize         float64     // <user set>   // Height of characters, set during loading (don't change after loading)
	Scale            float64     // = 1.f        // Base font scale, multiplied by the per-window font scale which you can adjust with SetFontScale()
	DisplayOffset    f64.Vec2    // = (0.f,1.f)  // Offset font rendering by xx pixels
	Glyphs           []FontGlyph //              // All glyphs.
	IndexAdvanceX    []float64   //              // Sparse. Glyphs->AdvanceX in a directly indexable way (more cache-friendly, for CalcTextSize functions which are often bottleneck in large UI).
	IndexLookup      []int       //              // Sparse. Index glyphs by Unicode code-point.
	FallbackGlyph    *FontGlyph  // == FindGlyph(FontFallbackChar)
	FallbackAdvanceX float64     // == FallbackGlyph->AdvanceX
	FallbackChar     rune        // = '?'        // Replacement glyph if one isn't found. Only set via SetFallbackChar()

	// Members: Cold ~18/26 bytes
	ConfigDataCount     int         // ~ 1          // Number of ImFontConfig involved in creating this font. Bigger than 1 when merging multiple font sources into one ImFont.
	ConfigData          *FontConfig //              // Pointer within ContainerAtlas->ConfigData
	ContainerAtlas      *FontAtlas  //              // What we has been loaded into
	Ascent, Descent     float64     //              // Ascent: distance from top to bottom of e.g. 'A' [0..FontSize]
	DirtyLookupTables   bool
	MetricsTotalSurface int //              // Total surface in pixels to get an idea of the font rasterization/texture cost (not exact, we approximate the cost of padding between glyphs)
}

type FontConfig struct {
	FontData                 []uint8  //          // TTF/OTF data
	FontDataSize             int      //          // TTF/OTF data size
	FontDataOwnedByAtlas     bool     // true     // TTF/OTF data ownership taken by the container ImFontAtlas (will delete memory itself).
	FontNo                   int      // 0        // Index of font within TTF/OTF file
	SizePixels               float64  //          // Size in pixels for rasterizer.
	OversampleH, OversampleV int      // 3, 1     // Rasterize at higher quality for sub-pixel positioning. We don't use sub-pixel positions on the Y axis.
	PixelSnapH               bool     // false    // Align every glyph to pixel boundary. Useful e.g. if you are merging a non-pixel aligned font with the default font. If enabled, you can set OversampleH/V to 1.
	GlyphExtraSpacing        f64.Vec2 // 0, 0     // Extra spacing (in pixels) between glyphs. Only X axis is supported for now.
	GlyphOffset              f64.Vec2 // 0, 0     // Offset all glyphs from this font input.
	GlyphRanges              []rune   // NULL     // Pointer to a user-provided list of Unicode range (2 value per range, values are inclusive, zero-terminated list). THE ARRAY DATA NEEDS TO PERSIST AS LONG AS THE FONT IS ALIVE.
	MergeMode                bool     // false    // Merge into previous ImFont, so you can combine multiple inputs font into one ImFont (e.g. ASCII font + icons + Japanese glyphs). You may want to use GlyphOffset.y when merge font of different heights.
	RasterizerFlags          uint     // 0x00     // Settings for custom font rasterizer (e.g. ImGuiFreeType). Leave as zero if you aren't using one.
	RasterizerMultiply       float64  // 1.0f     // Brighten (>1.0f) or darken (<1.0f) font output. Brightening small fonts may be a good workaround to make them more readable.

	// [Internal]
	Name    string // Name (strictly to ease debugging)
	DstFont *Font
}

type FontGlyph struct {
	Codepoint      rune    // 0x0000..0xFFFF
	AdvanceX       float64 // Distance to next character (= data from font + ImFontConfig::GlyphExtraSpacing.x baked in)
	X0, Y0, X1, Y1 float64 // Glyph corners
	U0, V0, U1, V1 float64 // Texture coordinates
}

func NewFont() *Font {
	f := &Font{}
	f.Init()
	return f
}

func (f *Font) Init() {
	f.Scale = 1
	f.FallbackChar = '?'
	f.DisplayOffset = f64.Vec2{0, 0}
	f.ClearOutputData()
}

func (f *Font) ClearOutputData() {
	f.FontSize = 0.0
	f.Glyphs = f.Glyphs[:0]
	f.IndexAdvanceX = f.IndexAdvanceX[:0]
	f.IndexLookup = f.IndexLookup[:0]
	f.FallbackGlyph = nil
	f.FallbackAdvanceX = 0.0
	f.ConfigDataCount = 0
	f.ConfigData = nil
	f.ContainerAtlas = nil
	f.Ascent = 0.0
	f.Descent = 0.0
	f.DirtyLookupTables = true
	f.MetricsTotalSurface = 0
}

func (c *Context) GetFont() *Font {
	return c.Font
}

func (c *Context) SetCurrentFont(font *Font) {
	c.Font = font
	c.FontBaseSize = c.IO.FontGlobalScale * c.Font.FontSize * c.Font.Scale
	c.FontSize = 0
	if c.CurrentWindow != nil {
		c.FontSize = c.CurrentWindow.CalcFontSize()
	}

	atlas := c.Font.ContainerAtlas
	c.DrawListSharedData.TexUvWhitePixel = atlas.TexUvWhitePixel
	c.DrawListSharedData.Font = c.Font
	c.DrawListSharedData.FontSize = c.FontSize
}

func (c *Context) PushFont(font *Font) {
	if font == nil {
		font = c.GetDefaultFont()
	}
	c.SetCurrentFont(font)
	c.FontStack = append(c.FontStack, font)
	c.CurrentWindow.DrawList.PushTextureID(font.ContainerAtlas.TexID)
}

func (c *Context) PopFont() {
	c.CurrentWindow.DrawList.PopTextureID()
	c.FontStack = c.FontStack[:len(c.FontStack)-1]
	if len(c.FontStack) == 0 {
		c.SetCurrentFont(c.GetDefaultFont())
	} else {
		c.SetCurrentFont(c.FontStack[len(c.FontStack)-1])
	}
}

func (c *Context) GetDefaultFont() *Font {
	if c.IO.FontDefault != nil {
		return c.IO.FontDefault
	}
	return c.IO.Fonts.Fonts[0]
}

func (c *Context) GetFontSize() float64 {
	return c.FontSize
}

func (c *Context) SetWindowFontScale(scale float64) {
	window := c.GetCurrentWindow()
	window.FontWindowScale = scale
	c.FontSize = window.CalcFontSize()
	c.DrawListSharedData.FontSize = c.FontSize
}

func (f *Font) CalcTextSizeA(size, max_width, wrap_width float64, text string) (text_size f64.Vec2, remaining int) {
	line_height := size
	scale := size / f.FontSize

	line_width := 0.0

	word_wrap_enabled := wrap_width > 0.0
	word_wrap_eol := -1

	s := 0
	for s < len(text) {
		if word_wrap_enabled {
			// Calculate how far we can render. Requires two passes on the string data but keeps the code simple and not intrusive for what's essentially an uncommon feature.
			if word_wrap_eol < 0 {
				word_wrap_eol = f.CalcWordWrapPositionA(scale, text[s:], wrap_width-line_width) + s
				// Wrap_width is too small to fit anything. Force displaying 1 character to minimize the height discontinuity.
				if word_wrap_eol == s {
					// +1 may not be a character start point in UTF-8 but it's ok because we use s >= word_wrap_eol below
					word_wrap_eol++
				}
			}

			if s >= word_wrap_eol {
				if text_size.X < line_width {
					text_size.X = line_width
				}
				text_size.Y += line_height
				line_width = 0
				word_wrap_eol = -1

				// Wrapping skips upcoming blanks
				for s < len(text) {
					c := rune(text[s])
					if CharIsSpace(c) {
						s++
					} else if c == '\n' {
						s++
						break
					} else {
						break
					}
				}
				continue
			}
		}

		prev_s := s
		c, size := utf8.DecodeRuneInString(text[s:])
		s += size

		switch c {
		case '\n':
			text_size.X = math.Max(text_size.X, line_width)
			text_size.Y += line_height
			line_width = 0.0
			continue
		case '\r':
			continue
		}

		char_width := f.FallbackAdvanceX * scale
		if int(c) < len(f.IndexAdvanceX) {
			char_width = f.IndexAdvanceX[c] * scale
		}
		if line_width+char_width >= max_width {
			s = prev_s
			break
		}

		line_width += char_width
	}

	if text_size.X < line_width {
		text_size.X = line_width
	}

	if line_width > 0 || text_size.Y == 0 {
		text_size.Y += line_height
	}

	remaining = s
	return
}

func (f *Font) CalcWordWrapPositionA(scale float64, text string, wrap_width float64) int {
	// Simple word-wrapping for English, not full-featured. Please submit failing cases!
	// FIXME: Much possible improvements (don't cut things like "word !", "word!!!" but cut within "word,,,,", more sensible support for punctuations, support for Unicode punctuations, etc.)

	// For references, possible wrap point marked with ^
	//  "aaa bbb, ccc,ddd. eee   fff. ggg!"
	//      ^    ^    ^   ^   ^__    ^    ^

	// List of hardcoded separators: .,;!?'"

	// Skip extra blanks after a line returns (that includes not counting them in width computation)
	// e.g. "Hello    world" --> "Hello" "World"

	// Cut words that cannot possibly fit within one line.
	// e.g.: "The tropical fish" with ~5 characters worth of width --> "The tr" "opical" "fish"

	line_width := 0.0
	word_width := 0.0
	blank_width := 0.0
	// We work with unscaled widths to avoid scaling every characters
	wrap_width /= scale

	word_end := 0
	prev_word_end := -1
	inside_word := true

	var s int
	for s < len(text) {
		c, size := utf8.DecodeRuneInString(text[s:])
		next_s := s + size

		switch c {
		case '\n':
			line_width, word_width, blank_width = 0, 0, 0
			inside_word = true
			s = next_s
			continue

		case '\r':
			s = next_s
			continue
		}

		char_width := f.FallbackAdvanceX
		if int(c) < len(f.IndexAdvanceX) {
			char_width = f.IndexAdvanceX[c]
		}
		if CharIsSpace(c) {
			if inside_word {
				line_width += blank_width
				blank_width = 0
				word_end = s
			}
			blank_width += char_width
			inside_word = false
		} else {
			word_width += char_width
			if inside_word {
				word_end = next_s
			} else {
				prev_word_end = word_end
				line_width += word_width + blank_width
				word_width, blank_width = 0, 0
			}

			// Allow wrapping after punctuation.
			inside_word = !(c == '.' || c == ',' || c == ';' || c == '!' || c == '?' || c == '"')
		}

		// We ignore blank width at the end of the line (they can be skipped)
		if line_width+word_width >= wrap_width {
			// Words that cannot possibly fit within an entire line will be cut anywhere.
			if word_width < wrap_width {
				if prev_word_end != -1 {
					s = prev_word_end
				} else {
					s = word_end
				}
			}
			break
		}

		s = next_s
	}

	return s
}

func CharIsSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == 0x3000
}

func (f *FontConfig) Init() {
	f.FontData = nil
	f.FontDataSize = 0
	f.FontDataOwnedByAtlas = true
	f.FontNo = 0
	f.SizePixels = 0.0
	f.OversampleH = 3
	f.OversampleV = 1
	f.PixelSnapH = false
	f.GlyphExtraSpacing = f64.Vec2{0, 0}
	f.GlyphOffset = f64.Vec2{0, 0}
	f.GlyphRanges = nil
	f.MergeMode = false
	f.RasterizerFlags = 0x00
	f.RasterizerMultiply = 1
	f.Name = ""
	f.DstFont = nil
}

func (f *Font) IsLoaded() bool {
	return f.ContainerAtlas != nil
}

func (f *Font) BuildLookupTable() {
	max_codepoint := 0
	for i := range f.Glyphs {
		max_codepoint = mathutil.Max(max_codepoint, int(f.Glyphs[i].Codepoint))
	}

	assert(len(f.Glyphs) < 0xFFFF) // -1 is reserved
	f.IndexAdvanceX = f.IndexAdvanceX[:0]
	f.IndexLookup = f.IndexLookup[:0]
	f.DirtyLookupTables = false
	f.GrowIndex(max_codepoint + 1)
	for i := range f.Glyphs {
		codepoint := f.Glyphs[i].Codepoint
		f.IndexAdvanceX[codepoint] = f.Glyphs[i].AdvanceX
		f.IndexLookup[codepoint] = i
	}

	// Create a glyph to handle TAB
	// FIXME: Needs proper TAB handling but it needs to be contextualized (or we could arbitrary say that each string starts at "column 0" ?)
	if f.FindGlyph(' ') != nil {
		// So we can call this function multiple times
		if f.Glyphs[len(f.Glyphs)-1].Codepoint != '\t' {
			f.Glyphs = append(f.Glyphs, FontGlyph{})
		}
		tab_glyph := &f.Glyphs[len(f.Glyphs)-1]
		*tab_glyph = *f.FindGlyph(' ')
		tab_glyph.Codepoint = '\t'
		tab_glyph.AdvanceX *= 4
		f.IndexAdvanceX[tab_glyph.Codepoint] = tab_glyph.AdvanceX
		f.IndexLookup[tab_glyph.Codepoint] = len(f.Glyphs) - 1
	}

	f.FallbackGlyph = f.FindGlyphNoFallback(f.FallbackChar)
	f.FallbackAdvanceX = 0.0
	if f.FallbackGlyph != nil {
		f.FallbackAdvanceX = f.FallbackGlyph.AdvanceX
	}
	for i := 0; i < max_codepoint+1; i++ {
		if f.IndexAdvanceX[i] < 0.0 {
			f.IndexAdvanceX[i] = f.FallbackAdvanceX
		}
	}
}

func (f *Font) SetFallbackChar(c rune) {
	f.FallbackChar = c
	f.BuildLookupTable()
}

func (f *Font) GrowIndex(new_size int) {
	assert(len(f.IndexAdvanceX) == len(f.IndexLookup))
	if new_size <= len(f.IndexLookup) {
		return
	}
	for i := len(f.IndexAdvanceX); i < new_size; i++ {
		f.IndexAdvanceX = append(f.IndexAdvanceX, -1)
	}
	for i := len(f.IndexLookup); i < new_size; i++ {
		f.IndexAdvanceX = append(f.IndexAdvanceX, 0xffff)
	}
}

func (f *Font) FindGlyph(c rune) *FontGlyph {
	if int(c) >= len(f.IndexLookup) {
		return f.FallbackGlyph
	}
	i := f.IndexLookup[c]
	if i == 0xffff {
		return f.FallbackGlyph
	}
	return &f.Glyphs[i]
}

func (f *Font) FindGlyphNoFallback(c rune) *FontGlyph {
	if int(c) >= len(f.IndexLookup) {
		return nil
	}
	i := f.IndexLookup[c]
	if i == 0xffff {
		return nil
	}
	return &f.Glyphs[i]
}

func (f *Font) BuildLookupTables() {
}

func (f *Font) AddGlyph(codepoint rune, x0, y0, x1, y1, u0, v0, u1, v1, advance_x float64) {
	var glyph FontGlyph
	glyph.Codepoint = codepoint
	glyph.X0 = x0
	glyph.Y0 = y0
	glyph.X1 = x1
	glyph.Y1 = y1
	glyph.U0 = u0
	glyph.V0 = v0
	glyph.U1 = u1
	glyph.V1 = v1
	glyph.AdvanceX = advance_x + f.ConfigData.GlyphExtraSpacing.X // Bake spacing into AdvanceX
	if f.ConfigData.PixelSnapH {
		glyph.AdvanceX = float64(int(glyph.AdvanceX + 0.5))
	}
	f.Glyphs = append(f.Glyphs, glyph)

	// Compute rough surface usage metrics (+1 to account for average padding, +0.99 to round)
	f.DirtyLookupTables = true
	f.MetricsTotalSurface += int((glyph.U1-glyph.U0)*float64(f.ContainerAtlas.TexWidth)+1.99) * (int)((glyph.V1-glyph.V0)*float64(f.ContainerAtlas.TexHeight)+1.99)
}

func (f *Font) AddRemapChar(dst, src rune, overwrite_dst bool) {
	// Currently this can only be called AFTER the font has been built, aka after calling ImFontAtlas::GetTexDataAs*() function.
	assert(len(f.IndexLookup) > 0)
	index_size := len(f.IndexLookup)

	// 'dst' already exists
	if int(dst) < index_size && f.IndexLookup[dst] == 0xFFFF && !overwrite_dst {
		return
	}

	// both 'dst' and 'src' don't exist -> no-op
	if int(src) >= index_size && int(dst) >= index_size {
		return
	}

	f.GrowIndex(int(dst) + 1)
	f.IndexLookup[dst] = 0xFFFF
	f.IndexAdvanceX[dst] = 1
	if int(src) < index_size {
		f.IndexLookup[dst] = f.IndexLookup[src]
		f.IndexAdvanceX[dst] = f.IndexAdvanceX[src]
	}
}