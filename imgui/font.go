package imgui

import "github.com/qeedquan/go-media/math/f64"

type CustomRect struct {
	ID            uint     // Input    // User ID. Use <0x10000 to map into a font glyph, >=0x10000 for other/internal/custom texture data.
	Width, Height uint     // Input    // Desired rectangle dimension
	X, Y          uint     // Output   // Packed position in Atlas
	GlyphAdvanceX float64  // Input    // For custom font glyphs only (ID<0x10000): glyph xadvance
	GlyphOffset   f64.Vec2 // Input    // For custom font glyphs only (ID<0x10000): glyph display offset
	Font          *Font    // Input    // For custom font glyphs only (ID<0x10000): target font
}

type FontAtlasFlags int

const (
	FontAtlasFlagsNoPowerOfTwoHeight FontAtlasFlags = 1 << 0 // Don't round the height to next power of two
	FontAtlasFlagsNoMouseCursors     FontAtlasFlags = 1 << 1 // Don't build software mouse cursors into the atlas
)

type FontAtlas struct {
	Flags           FontAtlasFlags // Build flags (see ImFontAtlasFlags_)
	TexID           TextureID      // User data to refer to the texture once it has been uploaded to user's graphic systems. It is passed back to you during rendering via the ImDrawCmd structure.
	TexDesiredWidth int            // Texture width desired by user before Build(). Must be a power-of-two. If have many glyphs your graphics API have texture size restrictions you may want to increase texture width to decrease height.
	TexGlyphPadding int            // Padding between glyphs within texture in pixels. Defaults to 1.

	// [Internal]
	// NB: Access texture data via GetTexData*() calls! Which will setup a default font for you.
	TexPixelsAlpha8 []uint8      // 1 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight
	TexPixelsRGBA32 []uint32     // 4 component per pixel, each component is unsigned 8-bit. Total size = TexWidth * TexHeight * 4
	TexWidth        int          // Texture width calculated during Build().
	TexHeight       int          // Texture height calculated during Build().
	TexUvScale      f64.Vec2     // = (1.0f/TexWidth, 1.0f/TexHeight)
	TexUvWhitePixel f64.Vec2     // Texture coordinates to a white pixel
	Fonts           []*Font      // Hold all the fonts returned by AddFont*. Fonts[0] is the default font upon calling ImGui::NewFrame(), use ImGui::PushFont()/PopFont() to change the current font.
	CustomRects     []CustomRect // Rectangles for packing custom texture data into the atlas.
	ConfigData      []FontConfig // Internal data
	CustomRectIds   [1]int       // Identifiers of custom texture rectangle used by ImFontAtlas/ImDrawList
}

type Font struct {
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
