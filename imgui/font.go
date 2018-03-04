package imgui

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/qeedquan/go-media/math/f64"
)

type Font struct {
	FontSize         float64
	Scale            float64
	DisplayOffset    f64.Vec2
	Glyphs           []FontGlyph
	FallbackAdvanceX float64
	IndexAdvanceX    []float64
	IndexLookup      []int
	FallbackGlyph    *FontGlyph
	FallbackChar     rune
	Ascent, Descent  float64
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
	window := c.CurrentWindow
	c.Font = font
	c.FontBaseSize = io.FontGlobalScale * c.Font.FontSize * c.Font.Scale
	c.FontSize = 0
	if window != nil {
		c.FontSize = window.CalcFontSize()
	}
}

func (f *Font) CalcTextSizeA(size, maxWidth, wrapWidth float64, text string) (textSize f64.Vec2, remaining int) {
	lineHeight := size
	scale := size / f.FontSize

	lineWidth := 0.0

	wordWrapEnabled := wrapWidth > 0
	wordWrapEOL := -1
	for idx := 0; idx < len(text); {
		ch, size := utf8.DecodeRuneInString(text[idx:])
		if wordWrapEnabled {
			// Calculate how far we can render. Requires two passes on the string data but
			// keeps the code simple and not intrusive for what's essentially an uncommon feature.
			if wordWrapEOL < 0 {
				wordWrapEOL = f.CalcWordWrapPositionA(scale, text[idx:], wrapWidth-lineWidth)

				// Wrap_width is too small to fit anything. Force displaying 1 character to minimize the height discontinuity.
				if wordWrapEOL == idx {
					wordWrapEOL++
				}
			}

			if idx >= wordWrapEOL {
				if textSize.X < lineWidth {
					textSize.X = lineWidth
				}
				textSize.Y += lineHeight
				lineWidth = 0
				wordWrapEOL = -1

				// wrapping skips upcoming blanks
				idx += skipSpace(text[idx:])
				continue
			}
		}

		prevIdx := idx
		idx += size

		if ch == '\n' {
			textSize.X = math.Max(textSize.X, lineWidth)
			textSize.Y += lineHeight
			lineWidth = 0
			continue
		}
		if ch == '\r' {
			continue
		}

		charWidth := f.FallbackAdvanceX
		if int(ch) < len(f.IndexAdvanceX) {
			charWidth = f.IndexAdvanceX[ch]
		}
		charWidth *= scale

		if lineWidth+charWidth >= maxWidth {
			remaining = prevIdx
			break
		}
		lineWidth += charWidth
	}

	if textSize.X < lineWidth {
		textSize.X = lineWidth
	}
	if lineWidth > 0 || textSize.Y == 0 {
		textSize.Y += lineHeight
	}
	return
}

func skipSpace(str string) int {
	for i, ch := range str {
		if !unicode.IsSpace(ch) {
			return i
		}
	}
	return len(str)
}

func (f *Font) CalcWordWrapPositionA(scale float64, text string, wrapWidth float64) int {
	// Simple word-wrapping for English, not full-featured.

	// For references, possible wrap point marked with ^
	//  "aaa bbb, ccc,ddd. eee   fff. ggg!"
	//      ^    ^    ^   ^   ^__    ^    ^

	// List of hardcoded separators: .,;!?'"

	// Skip extra blanks after a line returns (that includes not counting them in width computation)
	// e.g. "Hello    world" --> "Hello" "World"

	// Cut words that cannot possibly fit within one line.
	// e.g.: "The tropical fish" with ~5 characters worth of width --> "The tr" "opical" "fish"

	lineWidth := 0.0
	wordWidth := 0.0
	blankWidth := 0.0
	insideWord := true
	wordEnd := 0
	prevWordEnd := -1

	// work with unscaled widths to avoid scaling every character
	wrapWidth /= scale

	var (
		idx int
		ch  rune
	)
	for idx, ch = range text {
		if ch == '\n' {
			lineWidth, wordWidth, blankWidth = 0, 0, 0
			insideWord = true
		}
		if ch == '\r' {
			continue
		}

		charWidth := f.FallbackAdvanceX
		if int(ch) < len(f.IndexAdvanceX) {
			charWidth = f.IndexAdvanceX[ch]
		}

		if unicode.IsSpace(ch) {
			if insideWord {
				lineWidth += blankWidth
				blankWidth = 0
				wordEnd = idx
			}
			blankWidth += charWidth
			insideWord = false
		} else {
			wordWidth += charWidth
			if insideWord {
				wordEnd = idx + 1
			} else {
				prevWordEnd = wordEnd
				lineWidth += wordWidth + blankWidth
				wordWidth, blankWidth = 0, 0
			}

			// Allow wrapping after punctuation.
			insideWord = true
			if strings.IndexRune(".,;!?\"", ch) >= 0 {
				insideWord = false
			}
		}

		// We ignore blank width at the end of the line (they can be skipped)
		if lineWidth+wordWidth >= wrapWidth {
			if wordWidth < wrapWidth {
				idx = wordEnd
				if prevWordEnd >= 0 {
					idx = prevWordEnd
				}
			}
			break
		}
	}

	return idx
}
