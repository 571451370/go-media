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
	Layers [2][]*DrawList
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

func (c *Context) NewFrame() {
	// Load settings on first frame
	if !c.SettingsLoaded {
		c.SettingsLoaded = true
	}

	c.Time += c.IO.DeltaTime
	c.FrameCount += 1
	c.TooltipOverrideCount = 0
	c.WindowsActiveCount = 0

	c.SetCurrentFont(c.GetDefaultFont())
	c.DrawListSharedData.ClipRectFullscreen = f64.Vec4{0, 0, c.IO.DisplaySize.X, c.IO.DisplaySize.Y}
	c.DrawListSharedData.CurveTessellationTol = c.Style.CurveTessellationTol

	c.OverlayDrawList.Clear()
	c.OverlayDrawList.PushTextureID(c.IO.Fonts.TexID)
	c.OverlayDrawList.PushClipRectFullScreen()
	c.OverlayDrawList.Flags = 0
	if c.Style.AntiAliasedLines {
		c.OverlayDrawList.Flags |= DrawListFlagsAntiAliasedLines
	}
	if c.Style.AntiAliasedFill {
		c.OverlayDrawList.Flags |= DrawListFlagsAntiAliasedFill
	}

	// Mark rendering data as invalid to prevent user who may have a handle on it to use it
	c.DrawData.Clear()

	// Clear reference to active widget if the widget isn't alive anymore
	if c.HoveredIdPreviousFrame == 0 {
		c.HoveredIdTimer = 0
	}
	c.HoveredIdPreviousFrame = c.HoveredId
	c.HoveredId = 0
	c.HoveredIdAllowOverlap = false
	if !c.ActiveIdIsAlive && c.ActiveIdPreviousFrame == c.ActiveId && c.ActiveId != 0 {
		c.ClearActiveID()
	}
	if c.ActiveId != 0 {
		c.ActiveIdTimer += c.IO.DeltaTime
	}
	c.ActiveIdPreviousFrame = c.ActiveId
	c.ActiveIdIsAlive = false
	c.ActiveIdIsJustActivated = false
	if c.ScalarAsInputTextId != 0 && c.ActiveId != c.ScalarAsInputTextId {
		c.ScalarAsInputTextId = 0
	}

	// Elapse drag & drop payload
	if c.DragDropActive && c.DragDropPayload.DataFrameCount+1 < c.FrameCount {
		c.ClearDragDrop()
		for i := range c.DragDropPayloadBufHeap {
			c.DragDropPayloadBufHeap[i] = 0
		}
		for i := range c.DragDropPayloadBufLocal {
			c.DragDropPayloadBufLocal[i] = 0
		}
	}
	c.DragDropAcceptIdPrev = c.DragDropAcceptIdCurr
	c.DragDropAcceptIdCurr = 0
	c.DragDropAcceptIdCurrRectSurface = math.MaxFloat32

	// Update keyboard input state
	copy(c.IO.KeysDownDurationPrev[:], c.IO.KeysDownDuration[:])
	for i := range c.IO.KeysDown {
		c.IO.KeysDownDuration[i] = -1
		if c.IO.KeysDown[i] {
			if c.IO.KeysDownDuration[i] < 0 {
				c.IO.KeysDownDuration[i] = 0
			} else {
				c.IO.KeysDownDuration[i] = c.IO.KeysDownDuration[i] + c.IO.DeltaTime
			}
		}
	}

	// Update gamepad/keyboard directional navigation
	c.NavUpdate()

	// Update mouse input state
	// If mouse just appeared or disappeared (usually denoted by -FLT_MAX component, but in reality we test for -256000.0f) we cancel out movement in MouseDelta
}

func (c *Context) Render() {
	if c.FrameCountEnded != c.FrameCount {
		c.EndFrame()
	}
	c.FrameCountRendered = c.FrameCount

	// Gather windows to render
	c.IO.MetricsRenderVertices = 0
	c.IO.MetricsRenderIndices = 0
	c.IO.MetricsActiveWindows = 0
	c.DrawDataBuilder.Clear()

	var window_to_render_front_most *Window
	if c.NavWindowingTarget != nil && c.NavWindowingTarget.Flags&WindowFlagsNoBringToFrontOnFocus == 0 {
		window_to_render_front_most = c.NavWindowingTarget
	}

	for _, window := range c.Windows {
		if window.Active && window.HiddenFrames <= 0 && window.Flags&WindowFlagsChildWindow == 0 &&
			window != window_to_render_front_most {
			c.AddWindowToDrawDataSelectLayer(window)
		}
	}

	// NavWindowingTarget is always temporarily displayed as the front-most window
	if window_to_render_front_most != nil && window_to_render_front_most.Active &&
		window_to_render_front_most.HiddenFrames <= 0 {
		c.AddWindowToDrawDataSelectLayer(window_to_render_front_most)
	}
	c.DrawDataBuilder.FlattenIntoSingleLayer()

	// Draw software mouse cursor if requested
}

// This is normally called by Render(). You may want to call it directly if you want to avoid calling Render() but the gain will be very minimal.
func (c *Context) EndFrame() {
	// Don't process EndFrame() multiple times.
	if c.FrameCountEnded == c.FrameCount {
		return
	}

	// Notify OS when our Input Method Editor cursor has moved (e.g. CJK inputs using Microsoft IME)
	if c.IO.ImeSetInputScreenPosFn != nil && c.OsImePosRequest.DistanceSquared(c.OsImePosSet) > 0.0001 {
		c.IO.ImeSetInputScreenPosFn(int(c.OsImePosRequest.X), int(c.OsImePosRequest.Y))
		c.OsImePosSet = c.OsImePosRequest
	}

	// Hide implicit "Debug" window if it hasn't been used
	if c.CurrentWindow != nil && !c.CurrentWindow.WriteAccessed {
		c.CurrentWindow.Active = false
	}
	c.End()

	if c.ActiveId == 0 && c.HoveredId == 0 {
		// Unless we just made a window/popup appear
		if c.NavWindow == nil || !c.NavWindow.Appearing {
			// Click to focus window and start moving (after we're done with all our widgets)
			if c.IO.MouseClicked[0] {
				if c.HoveredRootWindow != nil {
					// Set ActiveId even if the _NoMove flag is set, without it dragging away from a window with _NoMove would activate hover on other windows.
					c.FocusWindow(c.HoveredWindow)
					c.SetActiveID(c.HoveredWindow.MoveId, c.HoveredWindow)
					c.NavDisableHighlight = true
					c.ActiveIdClickOffset = c.IO.MousePos.Sub(c.HoveredRootWindow.Pos)
					if c.HoveredWindow.Flags&WindowFlagsNoMove == 0 && c.HoveredRootWindow.Flags&WindowFlagsNoMove == 0 {
						c.MovingWindow = c.HoveredWindow
					}
				} else if c.NavWindow != nil && c.GetFrontMostModalRootWindow() == nil {
					// Clicking on void disable focus
					c.FocusWindow(nil)
				}
			}

			// With right mouse button we close popups without changing focus
			// (The left mouse button path calls FocusWindow which will lead NewFrame->ClosePopupsOverWindow to trigger)
			if c.IO.MouseClicked[1] {
				// Find the top-most window between HoveredWindow and the front most Modal Window.
				// This is where we can trim the popup stack.
				modal := c.GetFrontMostModalRootWindow()
				hovered_window_above_modal := false
				if modal == nil {
					hovered_window_above_modal = true
				}
				for i := len(c.Windows) - 1; i >= 0 && hovered_window_above_modal == false; i-- {
					window := c.Windows[i]
					if window == modal {
						break
					}
					if window == c.HoveredWindow {
						hovered_window_above_modal = true
					}
				}
				if hovered_window_above_modal {
					c.ClosePopupsOverWindow(c.HoveredWindow)
				} else {
					c.ClosePopupsOverWindow(modal)
				}
			}
		}
	}

	// Sort the window list so that all child windows are after their parent
	// We cannot do that on FocusWindow() because childs may not exist yet
	c.WindowsSortBuffer = c.WindowsSortBuffer[:0]
	for _, window := range c.Windows {
		if window.Active && window.Flags&WindowFlagsChildWindow != 0 {
			continue
		}
		c.AddWindowToSortedBuffer(&c.WindowsSortBuffer, window)
	}
	c.Windows, c.WindowsSortBuffer = c.WindowsSortBuffer, c.Windows

	// Clear Input data for next frame
	c.IO.MouseWheel = 0
	c.IO.MouseWheelH = 0
	for i := range c.IO.InputCharacters {
		c.IO.InputCharacters[i] = 0
	}
	for i := range c.IO.NavInputs {
		c.IO.NavInputs[i] = 0
	}

	c.FrameCountEnded = c.FrameCount
}

func (c *Context) End() {
	window := c.CurrentWindow
	if window.DC.ColumnsSet != nil {
		c.EndColumns()
	}
	// Inner window clip rectangle
	c.PopClipRect()

	// Stop logging
	// FIXME: add more options for scope of logging
	if window.Flags&WindowFlagsChildWindow == 0 {
		c.LogFinish()
	}

	// Pop from window stack
	c.CurrentWindowStack = c.CurrentWindowStack[:len(c.CurrentWindowStack)-1]
	if window.Flags&WindowFlagsPopup != 0 {
		c.CurrentPopupStack = c.CurrentPopupStack[:len(c.CurrentPopupStack)-1]
	}
	if len(c.CurrentWindowStack) == 0 {
		c.SetCurrentWindow(nil)
	} else {
		c.SetCurrentWindow(c.CurrentWindowStack[len(c.CurrentWindowStack)-1])
	}
}

func (c *Context) RenderNavHighlight(bb f64.Rectangle, id ID) {
}

func (c *Context) RenderFrame(p_min, p_max f64.Vec2, col color.RGBA) {
	c.RenderFrameEx(p_min, p_max, col, true, 0)
}

func (c *Context) RenderFrameEx(p_min, p_max f64.Vec2, col color.RGBA, border bool, rounding float64) {
}

func (c *Context) RenderArrow(pos f64.Vec2, dir Dir) {
}

func (c *Context) RenderTextClipped(pos_min, pos_max f64.Vec2, text string, text_size_if_known *f64.Vec2) {
	c.RenderTextClippedEx(pos_min, pos_max, text, text_size_if_known, f64.Vec2{0, 0}, nil)
}

func (c *Context) RenderTextClippedEx(pos_min, pos_max f64.Vec2, text string, text_size_if_known *f64.Vec2, align f64.Vec2, clip_rect *f64.Rectangle) {
}

func (d *DrawList) PathClear() {
	d._Path = d._Path[:0]
}

func (d *DrawList) PathLineTo(pos f64.Vec2) {
	d._Path = append(d._Path, pos)
}

func (d *DrawList) PathStroke(col color.RGBA, closed bool) {
	d.PathStrokeEx(col, closed, 1)
}

func (d *DrawList) PathStrokeEx(col color.RGBA, closed bool, thickness float64) {
	d.AddPolyline(d._Path, col, closed, thickness)
	d.PathClear()
}

func (d *DrawList) AddLine(a, b f64.Vec2, col color.RGBA) {
	d.AddLineEx(a, b, col, 1)
}

func (d *DrawList) AddLineEx(a, b f64.Vec2, col color.RGBA, thickness float64) {
	if col.A == 0 {
		return
	}
	half := f64.Vec2{0.5, 0.5}
	d.PathLineTo(a.Add(half))
	d.PathLineTo(b.Add(half))
	d.PathStrokeEx(col, false, thickness)
}

func (d *DrawList) AddPolyline(points []f64.Vec2, col color.RGBA, closed bool, thickness float64) {
}

func (d *DrawList) AddRect(p_min, p_max f64.Vec2, col color.RGBA, rounding float64) {
}

func (d *DrawList) AddRectFilled(p_min, p_max f64.Vec2, col color.RGBA, rounding float64) {
}

func (d *DrawList) AddImage(user_texture_id TextureID, a, b f64.Vec2) {
	d.AddImageEx(user_texture_id, a, b, f64.Vec2{0, 0}, f64.Vec2{1, 1}, color.RGBA{0xff, 0xff, 0xff, 0xff})
}

func (d *DrawList) AddImageEx(user_texture_id TextureID, a, b, uv_a, uv_b f64.Vec2, col color.RGBA) {
}

func (d *DrawList) Clear() {
	d.CmdBuffer = d.CmdBuffer[:0]
	d.IdxBuffer = d.IdxBuffer[:0]
	d.VtxBuffer = d.VtxBuffer[:0]
	d.Flags = DrawListFlagsAntiAliasedLines | DrawListFlagsAntiAliasedFill
	d._VtxCurrentIdx = 0
	d._VtxWritePtr = 0
	d._IdxWritePtr = 0
	d._ClipRectStack = d._ClipRectStack[:0]
	d._TextureIdStack = d._TextureIdStack[:0]
	d._Path = d._Path[:0]
	d._ChannelsCurrent = 0
	d._ChannelsCount = 1
	// NB: Do not clear channels so our allocations are re-used after the first frame.
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

func (d *DrawList) PushClipRectFullScreen() {
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

func (d *DrawList) ChannelsMerge() {
	// Note that we never use or rely on channels.Size because it is merely a buffer that we never shrink back to 0 to keep all sub-buffers ready for use.
	if d._ChannelsCount <= 1 {
		return
	}

	d.ChannelsSetCurrent(0)

	length := len(d.CmdBuffer)
	if length > 0 && d.CmdBuffer[length-1].ElemCount == 0 {
		d.CmdBuffer = d.CmdBuffer[:length-1]
	}

	new_cmd_buffer_count := 0
	new_idx_buffer_count := 0
	for i := 1; i < d._ChannelsCount; i++ {
		ch := &d._Channels[i]
		length := len(d.CmdBuffer)
		if length > 0 && ch.CmdBuffer[length-1].ElemCount == 0 {
			ch.CmdBuffer = ch.CmdBuffer[:length-1]
		}
		new_cmd_buffer_count += len(ch.CmdBuffer)
		new_idx_buffer_count += len(ch.IdxBuffer)
	}

	d.CmdBuffer = append(d.CmdBuffer, make([]DrawCmd, new_cmd_buffer_count)...)
	d.IdxBuffer = append(d.IdxBuffer, make([]DrawIdx, new_idx_buffer_count)...)
	cmd_write := len(d.CmdBuffer) - new_cmd_buffer_count
	d._IdxWritePtr = len(d.IdxBuffer) - new_idx_buffer_count
	for i := 1; i < d._ChannelsCount; i++ {
		ch := &d._Channels[i]
		if length := len(ch.CmdBuffer); length > 0 {
			copy(d.CmdBuffer[cmd_write:], ch.CmdBuffer[:])
			cmd_write += length
		}
		if length := len(ch.IdxBuffer); length > 0 {
			copy(d.IdxBuffer[d._IdxWritePtr:], ch.IdxBuffer[:])
			d._IdxWritePtr += length
		}
	}

	d.UpdateClipRect() // We call this instead of AddDrawCmd(), so that empty channels won't produce an extra draw call.
	d._ChannelsCount = 1
}

func (d *DrawDataBuilder) FlattenIntoSingleLayer() {
	for n := 1; n < len(d.Layers); n++ {
		d.Layers[0] = append(d.Layers[0], d.Layers[n]...)
		d.Layers[n] = d.Layers[n][:0]
	}
}

func (d *DrawDataBuilder) Clear() {
	for i := range d.Layers {
		d.Layers[i] = d.Layers[i][:0]
	}
}

func (c *Context) AddWindowToDrawData(out_render_list *[]*DrawList, window *Window) {
	c.AddDrawListToDrawData(out_render_list, window.DrawList)
	for i := 0; i < len(window.DC.ChildWindows); i++ {
		child := window.DC.ChildWindows[i]
		// clipped children may have been marked not active
		if child.Active && child.HiddenFrames <= 0 {
			c.AddWindowToDrawData(out_render_list, child)
		}
	}
}

func (c *Context) AddDrawListToDrawData(out_render_list *[]*DrawList, draw_list *DrawList) {
	if len(draw_list.CmdBuffer) == 0 {
		return
	}

	// Remove trailing command if unused
	last_cmd := &draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1]
	if last_cmd.ElemCount == 0 && last_cmd.UserCallback == nil {
		length := len(draw_list.CmdBuffer) - 1
		draw_list.CmdBuffer = draw_list.CmdBuffer[:length]
		if length == 0 {
			return
		}
	}

	*out_render_list = append(*out_render_list, draw_list)
}