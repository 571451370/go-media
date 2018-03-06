package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/math/f64"
)

type DrawListSharedData struct {
	TexUvWhitePixel      f64.Vec2 // UV of white pixel in the atlas
	Font                 *Font    // Current/default font (optional, for simplified AddText overload)
	FontSize             float64  // Current/default font size (optional, for simplified AddText overload)
	CurveTessellationTol float64
	ClipRectFullscreen   f64.Vec4 // Value for PushClipRectFullscreen()

	// Const data
	// FIXME: Bake rounded corners fill/borders in atlas
	CircleVtx12 [12]f64.Vec2
}

type DrawList struct {
	// This is what you have to render
	CmdBuffer []DrawCmd  // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []DrawIdx  // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []DrawVert // Vertex buffer.

	// [Internal, used while building lists]
	Flags           DrawListFlags       // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.
	_Data           *DrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	_OwnerName      string              // Pointer to owner window's name for debugging
	_VtxCurrentIdx  uint                // [Internal] == VtxBuffer.Size
	_VtxWritePtr    []DrawVert          // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_IdxWritePtr    []DrawIdx           // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_ClipRectStack  []f64.Vec4          // [Internal]
	_TextureIdStack []TextureID         // [Internal]
	_Path           []f64.Vec2          // [Internal] current path building                   _ChannelsCurrent int   // [Internal] current channel number (0)
	_ChannelsCount  int                 // [Internal] number of active channels (1+)
	_Channels       []DrawChannel       // [Internal] draw channels for columns API (not resized down so _ChannelsCount may be smaller than _Channels.Size)
}

type DrawCmd struct {
	ElemCount        uint        // Number of indices (multiple of 3) to be rendered as triangles. Vertices are stored in the callee ImDrawList's vtx_buffer[] array, indices in idx_buffer[].
	ClipRect         f64.Vec4    // Clipping rectangle (x1, y1, x2, y2)
	TextureId        TextureID   // User-provided texture ID. Set by user in ImfontAtlas::SetTexID() for fonts or passed to Image*() functions. Ignore if never using images or multiple fonts atlas.
	UserCallback     func()      // If != NULL, call the function instead of rendering the vertices. clip_rect and texture_id will be set normally.
	UserCallbackData interface{} // The draw callback code can access this.
}

type DrawIdx uint32

type DrawVert struct {
	Pos f64.Vec2
	UV  f64.Vec2
	Col f64.Vec2
}

type DrawChannel struct {
	CmdBuffer []DrawCmd
	IdxBuffer []DrawIdx
}

type DrawDataBuilder struct {
	Layers [2]*DrawList
}

type DrawData struct {
	Valid         bool // Only valid after Render() is called and before the next NewFrame() is called.
	CmdLists      []*DrawList
	CmdListsCount int
	TotalVtxCount int // For convenience, sum of all cmd_lists vtx_buffer.Size
	TotalIdxCount int // For convenience, sum of all cmd_lists idx_buffer.Size
}

type DrawListFlags int

const (
	DrawListFlagsAntiAliasedLines DrawListFlags = 1 << 0
	DrawListFlagsAntiAliasedFill  DrawListFlags = 1 << 1
)

func (c *Context) RenderNavHighlight(bb f64.Rectangle, id ID) {
}

func (c *Context) RenderFrame(p_min, p_max f64.Vec2, col color.RGBA) {
	c.RenderFrameDx(p_min, p_max, col, true, 0)
}

func (c *Context) RenderFrameDx(p_min, p_max f64.Vec2, col color.RGBA, border bool, rounding float64) {
}

func (c *Context) RenderArrow(pos f64.Vec2, dir Dir) {
}
