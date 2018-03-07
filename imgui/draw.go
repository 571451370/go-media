package imgui

import (
	"image/color"
	"math"

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
	Flags            DrawListFlags       // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.
	_Data            *DrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	_OwnerName       string              // Pointer to owner window's name for debugging
	_VtxCurrentIdx   uint                // [Internal] == VtxBuffer.Size
	_VtxWritePtr     int                 // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_IdxWritePtr     int                 // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_ClipRectStack   []f64.Vec4          // [Internal]
	_TextureIdStack  []TextureID         // [Internal]
	_Path            []f64.Vec2          // [Internal] current path building                   _ChannelsCurrent int   // [Internal] current channel number (0)
	_ChannelsCurrent int                 // [Internal] current channel number (0)
	_ChannelsCount   int                 // [Internal] number of active channels (1+)
	_Channels        []DrawChannel       // [Internal] draw channels for columns API (not resized down so _ChannelsCount may be smaller than _Channels.Size)
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

func (d *DrawList) PopClipRect() {
	d._ClipRectStack = d._ClipRectStack[:len(d._ClipRectStack)-1]
	d.UpdateClipRect()
}

func (d *DrawList) PushTextureID(texture_id TextureID) {
	d._TextureIdStack = append(d._TextureIdStack, texture_id)
	d.UpdateTextureID()
}

func (d *DrawList) PopTextureID() {
	d._TextureIdStack = d._TextureIdStack[:len(d._TextureIdStack)-1]
	d.UpdateTextureID()
}

func (d *DrawList) UpdateTextureID() {
	// If current command is used with different settings we need to add a new command
	curr_texture_id := d.GetCurrentTextureId()
	var curr_cmd *DrawCmd
	if length := len(d.CmdBuffer); length > 0 {
		curr_cmd = &d.CmdBuffer[length-1]
	}
	if curr_cmd == nil || (curr_cmd.ElemCount != 0 && curr_cmd.TextureId == curr_texture_id) || curr_cmd.UserCallback != nil {
		d.AddDrawCmd()
		return
	}

	// Try to merge with previous command if it matches, else use current command
	var prev_cmd *DrawCmd
	if length := len(d.CmdBuffer); length > 1 {
		prev_cmd = &d.CmdBuffer[length-2]
	}
	if curr_cmd.ElemCount == 0 && prev_cmd != nil && prev_cmd.TextureId == curr_texture_id &&
		prev_cmd.ClipRect == d.GetCurrentClipRect() && prev_cmd.UserCallback == nil {
		d.CmdBuffer = d.CmdBuffer[:len(d.CmdBuffer)-1]
	} else {
		curr_cmd.TextureId = curr_texture_id
	}
}

func (d *DrawList) ChannelsSetCurrent(idx int) {
	if d._ChannelsCurrent == idx {
		return
	}
	d._Channels[d._ChannelsCurrent].CmdBuffer = d.CmdBuffer
	d._Channels[d._ChannelsCurrent].IdxBuffer = d.IdxBuffer

	d._ChannelsCurrent = idx

	d.CmdBuffer = d._Channels[d._ChannelsCurrent].CmdBuffer
	d.IdxBuffer = d._Channels[d._ChannelsCurrent].IdxBuffer
	d._IdxWritePtr = len(d.IdxBuffer)
}

func (d *DrawList) PushClipRect(cr_min, cr_max f64.Vec2) {
	d.PushClipRectEx(cr_min, cr_max, false)
}

func (d *DrawList) PushClipRectEx(cr_min, cr_max f64.Vec2, intersect_with_current_clip_rect bool) {
	cr := f64.Vec4{cr_min.X, cr_min.Y, cr_max.X, cr_max.Y}
	length := len(d._ClipRectStack)
	if intersect_with_current_clip_rect && length > 0 {
		current := d._ClipRectStack[length-1]
		if cr.X < current.X {
			cr.X = current.X
		}
		if cr.Y < current.Y {
			cr.Y = current.Y
		}
		if cr.Z > current.Z {
			cr.Z = current.Z
		}
		if cr.W > current.W {
			cr.W = current.W
		}
	}
	cr.Z = math.Max(cr.X, cr.Z)
	cr.W = math.Max(cr.Y, cr.W)

	d._ClipRectStack = append(d._ClipRectStack, cr)
	d.UpdateClipRect()
}

// Our scheme may appears a bit unusual, basically we want the most-common calls AddLine AddRect etc. to not have to perform any check so we always have a command ready in the stack.
func (d *DrawList) UpdateClipRect() {
	// If current command is used with different settings we need to add a new command
	curr_clip_rect := d.GetCurrentClipRect()
	var curr_cmd *DrawCmd
	if length := len(d.CmdBuffer); length > 0 {
		curr_cmd = &d.CmdBuffer[length-1]
	}
	if curr_cmd == nil || (curr_cmd.ElemCount != 0 && curr_cmd.ClipRect == curr_clip_rect) || curr_cmd.UserCallback != nil {
		d.AddDrawCmd()
		return
	}

	// Try to merge with previous command if it matches, else use current command
	var prev_cmd *DrawCmd
	if length := len(d.CmdBuffer); length > 1 {
		prev_cmd = &d.CmdBuffer[length-2]
	}

	if curr_cmd.ElemCount == 0 && prev_cmd != nil && prev_cmd.ClipRect == curr_clip_rect &&
		prev_cmd.TextureId == d.GetCurrentTextureId() && prev_cmd.UserCallback == nil {
		d.CmdBuffer = d.CmdBuffer[:len(d.CmdBuffer)-1]
	} else {
		curr_cmd.ClipRect = curr_clip_rect
	}
}

func (d *DrawList) PushClipRectFullscreen() {
	clipRect := d._Data.ClipRectFullscreen
	d.PushClipRect(f64.Vec2{clipRect.X, clipRect.Y}, f64.Vec2{clipRect.Z, clipRect.W})
}

func (d *DrawList) GetCurrentClipRect() f64.Vec4 {
	length := len(d._ClipRectStack)
	if length > 0 {
		return d._ClipRectStack[length-1]
	}
	return d._Data.ClipRectFullscreen
}

func (d *DrawList) GetCurrentTextureId() TextureID {
	length := len(d._TextureIdStack)
	if length > 0 {
		return d._TextureIdStack[length-1]
	}
	return nil
}

func (d *DrawList) AddDrawCmd() {
	var draw_cmd DrawCmd
	draw_cmd.ClipRect = d.GetCurrentClipRect()
	draw_cmd.TextureId = d.GetCurrentTextureId()
	d.CmdBuffer = append(d.CmdBuffer, draw_cmd)
}