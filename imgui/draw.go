package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type DrawCornerFlags int

const (
	DrawCornerFlagsTopLeft  DrawCornerFlags = 1 << 0                                            // 0x1
	DrawCornerFlagsTopRight DrawCornerFlags = 1 << 1                                            // 0x2
	DrawCornerFlagsBotLeft  DrawCornerFlags = 1 << 2                                            // 0x4
	DrawCornerFlagsBotRight DrawCornerFlags = 1 << 3                                            // 0x8
	DrawCornerFlagsTop      DrawCornerFlags = DrawCornerFlagsTopLeft | DrawCornerFlagsTopRight  // 0x3
	DrawCornerFlagsBot      DrawCornerFlags = DrawCornerFlagsBotLeft | DrawCornerFlagsBotRight  // 0xC
	DrawCornerFlagsLeft     DrawCornerFlags = DrawCornerFlagsTopLeft | DrawCornerFlagsBotLeft   // 0x5
	DrawCornerFlagsRight    DrawCornerFlags = DrawCornerFlagsTopRight | DrawCornerFlagsBotRight // 0xA
	DrawCornerFlagsAll      DrawCornerFlags = 0xF                                               // In your function calls you may use ~0 (= all bits sets) instead of DrawCornerFlags_All, as a convenience
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
	if c.IsMousePosValid(&c.IO.MousePos) && c.IsMousePosValid(&c.IO.MousePosPrev) {
		c.IO.MouseDelta = c.IO.MousePos.Sub(c.IO.MousePosPrev)
	} else {
		c.IO.MouseDelta = f64.Vec2{0, 0}
	}
	if c.IO.MouseDelta.X != 0 || c.IO.MouseDelta.Y != 0 {
		c.NavDisableMouseHover = false
	}

	c.IO.MousePosPrev = c.IO.MousePos
	for i := range c.IO.MouseDown {
		c.IO.MouseClicked[i] = c.IO.MouseDown[i] && c.IO.MouseDownDuration[i] < 0
		c.IO.MouseReleased[i] = !c.IO.MouseDown[i] && c.IO.MouseDownDuration[i] >= 0
		c.IO.MouseDownDurationPrev[i] = c.IO.MouseDownDuration[i]
		if c.IO.MouseDown[i] {
			if c.IO.MouseDownDuration[i] < 0 {
				c.IO.MouseDownDuration[i] = 0
			} else {
				c.IO.MouseDownDuration[i] = c.IO.MouseDownDuration[i] + c.IO.DeltaTime
			}
		} else {
			c.IO.MouseDownDuration[i] = -1
		}
		c.IO.MouseDoubleClicked[i] = false

		if c.IO.MouseClicked[i] {
			if c.Time-c.IO.MouseClickedTime[i] < c.IO.MouseDoubleClickTime {
				if c.IO.MousePos.DistanceSquared(c.IO.MouseClickedPos[i]) < c.IO.MouseDoubleClickMaxDist*c.IO.MouseDoubleClickMaxDist {
					c.IO.MouseDoubleClicked[i] = true
				}
				c.IO.MouseClickedTime[i] = -math.MaxFloat32 // so the third click isn't turned into a double-click
			} else {
				c.IO.MouseClickedTime[i] = c.Time
			}

			c.IO.MouseClickedPos[i] = c.IO.MousePos
			c.IO.MouseDragMaxDistanceAbs[i] = f64.Vec2{0, 0}
			c.IO.MouseDragMaxDistanceSqr[i] = 0
		} else if c.IO.MouseDown[i] {
			mouse_delta := c.IO.MousePos.Sub(c.IO.MouseClickedPos[i])
			c.IO.MouseDragMaxDistanceAbs[i].X = math.Max(c.IO.MouseDragMaxDistanceAbs[i].X, math.Abs(mouse_delta.X))
			c.IO.MouseDragMaxDistanceAbs[i].Y = math.Max(c.IO.MouseDragMaxDistanceAbs[i].Y, math.Abs(mouse_delta.Y))
			c.IO.MouseDragMaxDistanceSqr[i] = math.Max(c.IO.MouseDragMaxDistanceSqr[i], mouse_delta.LenSquared())
		}
		// Clicking any mouse button reactivate mouse hovering which may have been deactivated by gamepad/keyboard navigation
		if c.IO.MouseClicked[i] {
			c.NavDisableMouseHover = false
		}
	}

	// Calculate frame-rate for the user, as a purely luxurious feature
	c.FramerateSecPerFrameAccum += c.IO.DeltaTime - c.FramerateSecPerFrame[c.FramerateSecPerFrameIdx]
	c.FramerateSecPerFrame[c.FramerateSecPerFrameIdx] = c.IO.DeltaTime
	c.FramerateSecPerFrameIdx = (c.FramerateSecPerFrameIdx + 1) % len(c.FramerateSecPerFrame)
	c.IO.Framerate = 1.0 / (c.FramerateSecPerFrameAccum / float64(len(c.FramerateSecPerFrame)))

	// Handle user moving window with mouse (at the beginning of the frame to avoid input lag or sheering)
	c.UpdateMovingWindow()

	// Delay saving settings so we don't spam disk too much
	if c.SettingsDirtyTimer > 0 {
		c.SettingsDirtyTimer -= c.IO.DeltaTime
		if c.SettingsDirtyTimer <= 0 {
			c.SaveIniSettingsToDisk(c.IO.IniFilename)
		}
	}

	// Find the window we are hovering
	// - Child windows can extend beyond the limit of their parent so we need to derive HoveredRootWindow from HoveredWindow.
	// - When moving a window we can skip the search, which also conveniently bypasses the fact that window->WindowRectClipped is lagging as this point.
	// - We also support the moved window toggling the NoInputs flag after moving has started in order to be able to detect windows below it, which is useful for e.g. docking mechanisms.
	if c.MovingWindow != nil && c.MovingWindow.Flags&WindowFlagsNoInputs == 0 {
		c.HoveredWindow = c.MovingWindow
	} else {
		c.HoveredWindow = c.FindHoveredWindow()
	}
	c.HoveredRootWindow = nil
	if c.HoveredWindow != nil {
		c.HoveredRootWindow = c.HoveredWindow.RootWindow
	}

	modal_window := c.GetFrontMostModalRootWindow()
	if modal_window != nil {
		c.ModalWindowDarkeningRatio = math.Min(c.ModalWindowDarkeningRatio+c.IO.DeltaTime*6, 1)
		if c.HoveredRootWindow != nil && !c.IsWindowChildOf(c.HoveredRootWindow, modal_window) {
			c.HoveredRootWindow = nil
			c.HoveredWindow = nil
		}
	} else {
		c.ModalWindowDarkeningRatio = 0
	}

	// Update the WantCaptureMouse/WantCaptureKeyboard flags, so user can capture/discard the inputs away from the rest of their application.
	// When clicking outside of a window we assume the click is owned by the application and won't request capture. We need to track click ownership.
	mouse_earliest_button_down := -1
	mouse_any_down := false
	for i := range c.IO.MouseDown {
		if c.IO.MouseClicked[i] {
			c.IO.MouseDownOwned[i] = c.HoveredWindow != nil || len(c.OpenPopupStack) > 0
		}
		if c.IO.MouseDown[i] {
			mouse_any_down = true
		}
		if c.IO.MouseDown[i] {
			if mouse_earliest_button_down == -1 || c.IO.MouseClickedTime[i] < c.IO.MouseClickedTime[mouse_earliest_button_down] {
				mouse_earliest_button_down = i
			}
		}
	}
	mouse_avail_to_imgui := (mouse_earliest_button_down == -1) || c.IO.MouseDownOwned[mouse_earliest_button_down]
	if c.WantCaptureMouseNextFrame != -1 {
		c.IO.WantCaptureMouse = c.WantCaptureMouseNextFrame != 0
	} else {
		c.IO.WantCaptureMouse = (mouse_avail_to_imgui && (c.HoveredWindow != nil || mouse_any_down)) || len(c.OpenPopupStack) > 0
	}

	if c.WantCaptureKeyboardNextFrame != -1 {
		c.IO.WantCaptureKeyboard = c.WantCaptureKeyboardNextFrame != 0
	} else {
		c.IO.WantCaptureKeyboard = c.ActiveId != 0 || modal_window != nil
	}
	if c.IO.NavActive && c.IO.ConfigFlags&ConfigFlagsNavEnableKeyboard != 0 && c.IO.ConfigFlags&ConfigFlagsNavNoCaptureKeyboard == 0 {
		c.IO.WantCaptureKeyboard = true
	}

	c.IO.WantTextInput = false
	if c.WantTextInputNextFrame != -1 {
		c.IO.WantTextInput = c.WantTextInputNextFrame != 0
	}
	c.MouseCursor = MouseCursorArrow
	c.WantCaptureMouseNextFrame = -1
	c.WantCaptureKeyboardNextFrame = -1
	c.WantTextInputNextFrame = -1
	c.OsImePosRequest = f64.Vec2{1, 1} // OS Input Method Editor showing on top-left of our window by default

	// If mouse was first clicked outside of ImGui bounds we also cancel out hovering.
	// FIXME: For patterns of drag and drop across OS windows, we may need to rework/remove this test (first committed 311c0ca9 on 2015/02)
	mouse_dragging_extern_payload := c.DragDropActive && c.DragDropSourceFlags&DragDropFlagsSourceExtern != 0
	if !mouse_avail_to_imgui && !mouse_dragging_extern_payload {
		c.HoveredWindow = nil
		c.HoveredRootWindow = nil
	}

	// Mouse wheel scrolling, scale
	if c.HoveredWindow != nil && !c.HoveredWindow.Collapsed && (c.IO.MouseWheel != 0 || c.IO.MouseWheelH != 0) {
		// If a child window has the ImGuiWindowFlags_NoScrollWithMouse flag, we give a chance to scroll its parent (unless either ImGuiWindowFlags_NoInputs or ImGuiWindowFlags_NoScrollbar are also set).
		window := c.HoveredWindow
		scroll_window := window
		for scroll_window.Flags&WindowFlagsChildWindow != 0 &&
			scroll_window.Flags&WindowFlagsNoScrollWithMouse != 0 &&
			scroll_window.Flags&WindowFlagsNoScrollbar == 0 &&
			scroll_window.Flags&WindowFlagsNoInputs == 0 &&
			scroll_window.ParentWindow != nil {
			scroll_window = scroll_window.ParentWindow
		}
		scroll_allowed := scroll_window.Flags&WindowFlagsNoScrollWithMouse == 0 && scroll_window.Flags&WindowFlagsNoInputs == 0

		if c.IO.MouseWheel != 0 {
			if c.IO.KeyCtrl && c.IO.FontAllowUserScaling {
				// Zoom / Scale window
				new_font_scale := f64.Clamp(window.FontWindowScale+c.IO.MouseWheel*0.10, 0.50, 2.50)
				scale := new_font_scale / window.FontWindowScale
				window.FontWindowScale = new_font_scale

				offset := f64.Vec2{
					window.Size.X * (1.0 - scale) * (c.IO.MousePos.X - window.Pos.X) / window.Size.X,
					window.Size.Y * (1.0 - scale) * (c.IO.MousePos.Y - window.Pos.Y) / window.Size.Y,
				}
				window.Pos = window.Pos.Add(offset)
				window.PosFloat = window.PosFloat.Add(offset)
				window.Size = window.Size.Scale(scale)
				window.SizeFull = window.SizeFull.Scale(scale)
			} else if !c.IO.KeyCtrl && scroll_allowed {
				// Mouse wheel vertical scrolling
				scroll_amount := 5 * scroll_window.CalcFontSize()
				scroll_amount = math.Min(
					scroll_amount,
					(scroll_window.ContentsRegionRect.Dy()+scroll_window.WindowPadding.Y*2.0)*0.67,
				)
				c.SetWindowScrollY(scroll_window, scroll_window.Scroll.Y-c.IO.MouseWheel*scroll_amount)
			}
		}
		if c.IO.MouseWheelH != 0 && scroll_allowed {
			// Mouse wheel horizontal scrolling (for hardware that supports it)
			scroll_amount := scroll_window.CalcFontSize()
			if !c.IO.KeyCtrl && window.Flags&WindowFlagsNoScrollWithMouse == 0 {
				c.SetWindowScrollX(window, window.Scroll.X-c.IO.MouseWheelH*scroll_amount)
			}
		}
	}

	// Pressing TAB activate widget focus
	if c.ActiveId == 0 && c.NavWindow != nil && c.NavWindow.Active && c.NavWindow.Flags&WindowFlagsNoNavInputs == 0 &&
		!c.IO.KeyCtrl && c.IsKeyPressedMap(KeyTab, false) {
		if c.NavId != 0 && c.NavIdTabCounter != math.MaxInt32 {
			c.NavWindow.FocusIdxTabRequestNext = c.NavIdTabCounter + 1
			if c.IO.KeyShift {
				c.NavWindow.FocusIdxTabRequestNext -= 1
			} else {
				c.NavWindow.FocusIdxTabRequestNext += 1
			}
		} else {
			c.NavWindow.FocusIdxTabRequestNext = 0
			if c.IO.KeyShift {
				c.NavWindow.FocusIdxTabRequestNext = -1
			}
		}
	}
	c.NavIdTabCounter = math.MaxInt32

	// Mark all windows as not visible
	for i := range c.Windows {
		window := c.Windows[i]
		window.WasActive = window.Active
		window.Active = false
		window.WriteAccessed = false
	}

	// Closing the focused window restore focus to the first active root window in descending z-order
	if c.NavWindow != nil && !c.NavWindow.WasActive {
		c.FocusFrontMostActiveWindow(nil)
	}

	// No window should be open at the beginning of the frame.
	// But in order to allow the user to call NewFrame() multiple times without calling Render(), we are doing an explicit clear.
	c.CurrentWindowStack = c.CurrentWindowStack[:0]
	c.CurrentPopupStack = c.CurrentPopupStack[:0]
	c.ClosePopupsOverWindow(c.NavWindow)

	// Create implicit window - we will only render it if the user has added something to it.
	// We don't use "Debug" to avoid colliding with user trying to create a "Debug" window with custom flags.
	c.SetNextWindowSize(f64.Vec2{400, 400}, CondFirstUseEver)
	c.Begin("Debug##Default")
}

func (c *Context) Begin(name string) bool {
	return c.BeginEx(name, nil, 0)
}

// Push a new ImGui window to add widgets to.
// - A default window called "Debug" is automatically stacked at the beginning of every frame so you can use widgets without explicitly calling a Begin/End pair.
// - Begin/End can be called multiple times during the frame with the same window name to append content.
// - The window name is used as a unique identifier to preserve window information across frames (and save rudimentary information to the .ini file).
//   You can use the "##" or "###" markers to use the same label with different id, or same id with different label. See documentation at the top of this file.
// - Return false when window is collapsed, so you can early out in your code. You always need to call ImGui::End() even if false is returned.
// - Passing 'bool* p_open' displays a Close button on the upper-right corner of the window, the pointed value will be set to false when the button is pressed.
func (c *Context) BeginEx(name string, p_open *bool, flags WindowFlags) bool {
	// Find or create
	style := &c.Style
	window := c.FindWindowByName(name)
	if window == nil {
		// Any condition flag will do since we are creating a new window here.
		var size_on_first_use f64.Vec2
		if c.NextWindowData.SizeCond != 0 {
			size_on_first_use = c.NextWindowData.SizeVal
			window = c.CreateNewWindow(name, size_on_first_use, flags)
		}
	}

	// Automatically disable manual moving/resizing when NoInputs is set
	if flags&WindowFlagsNoInputs != 0 {
		flags |= WindowFlagsNoMove | WindowFlagsNoResize
	}

	current_frame := c.FrameCount
	first_begin_of_the_frame := window.LastFrameActive != current_frame
	if first_begin_of_the_frame {
		window.Flags = flags
	} else {
		flags = window.Flags
	}

	// Update the Appearing flag
	// Not using !WasActive because the implicit "Debug" window would always toggle off->on
	window_just_activated_by_user := window.LastFrameActive < current_frame-1
	window_just_appearing_after_hidden_for_resize := window.HiddenFrames == 1
	if flags&WindowFlagsPopup != 0 {
		popup_ref := c.OpenPopupStack[len(c.CurrentPopupStack)]
		// We recycle popups so treat window as activated if popup id changed
		if window.PopupId != popup_ref.PopupId {
			window_just_activated_by_user = true
		}
		if window != popup_ref.Window {
			window_just_activated_by_user = true
		}
	}
	window.Appearing = window_just_activated_by_user || window_just_appearing_after_hidden_for_resize
	window.CloseButton = p_open != nil
	if window.Appearing {
		c.SetWindowConditionAllowFlags(window, CondAppearing, true)
	}

	// Parent window is latched only on the first call to Begin() of the frame, so further append-calls can be done from a different window stack
	var parent_window_in_stack *Window
	if len(c.CurrentWindowStack) > 0 {
		parent_window_in_stack = c.CurrentWindowStack[len(c.CurrentWindowStack)-1]
	}

	// Add to stack
	c.CurrentWindowStack = append(c.CurrentWindowStack, window)
	c.SetCurrentWindow(window)
	if flags&WindowFlagsPopup != 0 {
		popup_ref := c.OpenPopupStack[len(c.CurrentPopupStack)]
		popup_ref.Window = window
		c.CurrentPopupStack = append(c.CurrentPopupStack, popup_ref)
		window.PopupId = popup_ref.PopupId
	}

	if window_just_appearing_after_hidden_for_resize && flags&WindowFlagsChildWindow == 0 {
		window.NavLastIds[0] = 0
	}

	// Process SetNextWindow***() calls
	window_pos_set_by_api := false
	window_size_x_set_by_api := false
	window_size_y_set_by_api := false
	if c.NextWindowData.PosCond != 0 {
		window_pos_set_by_api = window.SetWindowPosAllowFlags&c.NextWindowData.PosCond != 0
		if window_pos_set_by_api && c.NextWindowData.PosPivotVal.LenSquared() > 0.00001 {
			// May be processed on the next frame if this is our first frame and we are measuring size
			// FIXME: Look into removing the branch so everything can go through this same code path for consistency.
			window.SetWindowPosVal = c.NextWindowData.PosVal
			window.SetWindowPosPivot = c.NextWindowData.PosPivotVal
			window.SetWindowPosAllowFlags &^= (CondOnce | CondFirstUseEver | CondAppearing)
		} else {
			c.SetWindowPos(window, c.NextWindowData.PosVal, c.NextWindowData.PosCond)
		}
		c.NextWindowData.PosCond = 0
	}

	if c.NextWindowData.SizeCond != 0 {
		window_size_x_set_by_api = (window.SetWindowSizeAllowFlags&c.NextWindowData.SizeCond) != 0 && (c.NextWindowData.SizeVal.X > 0.0)
		window_size_y_set_by_api = (window.SetWindowSizeAllowFlags&c.NextWindowData.SizeCond) != 0 && (c.NextWindowData.SizeVal.Y > 0.0)
		c.SetWindowSize(window, c.NextWindowData.SizeVal, c.NextWindowData.SizeCond)
		c.NextWindowData.SizeCond = 0
	}

	if c.NextWindowData.ContentSizeCond != 0 {
		// Adjust passed "client size" to become a "window size"
		window.SizeContentsExplicit = c.NextWindowData.ContentSizeVal
		if window.SizeContentsExplicit.Y != 0.0 {
			window.SizeContentsExplicit.Y += window.TitleBarHeight() + window.MenuBarHeight()
		}
		c.NextWindowData.ContentSizeCond = 0
	} else if first_begin_of_the_frame {
		window.SizeContentsExplicit = f64.Vec2{0, 0}
	}

	if c.NextWindowData.CollapsedCond != 0 {
		c.SetWindowCollapsed(window, c.NextWindowData.CollapsedVal, c.NextWindowData.CollapsedCond)
		c.NextWindowData.CollapsedCond = 0

	}

	if c.NextWindowData.FocusCond != 0 {
		c.SetWindowFocus()
		c.NextWindowData.FocusCond = 0
	}

	if window.Appearing {
		c.SetWindowConditionAllowFlags(window, CondAppearing, false)
	}

	var parent_window *Window
	// When reusing window again multiple times a frame, just append content (don't need to setup again)
	if first_begin_of_the_frame {
		// FIXME-WIP: Undocumented behavior of Child+Tooltip for pinned tooltip (#1345)
		window_is_child_tooltip := (flags&WindowFlagsChildWindow) != 0 && (flags&WindowFlagsTooltip) != 0

		// Initialize
		window.ParentWindow = parent_window
		window.RootWindow = window
		window.RootWindowForTitleBarHighlight = window
		window.RootWindowForTabbing = window
		window.RootWindowForNav = window
		if parent_window != nil && flags&WindowFlagsChildWindow != 0 && !window_is_child_tooltip {
			window.RootWindow = parent_window.RootWindow
		}

		if parent_window != nil && flags&WindowFlagsModal == 0 && flags&(WindowFlagsChildWindow|WindowFlagsPopup) != 0 {
			// Same value in master branch, will differ for docking
			window.RootWindowForTitleBarHighlight = parent_window.RootWindowForTitleBarHighlight
			window.RootWindowForTabbing = window.RootWindowForTitleBarHighlight
		}

		for window.RootWindowForNav.Flags&WindowFlagsNavFlattened != 0 {
			window.RootWindowForNav = window.RootWindowForNav.ParentWindow
		}

		window.Active = true
		window.BeginOrderWithinParent = 0
		window.BeginOrderWithinContext = c.WindowsActiveCount
		c.WindowsActiveCount++
		window.BeginCount = 0
		window.ClipRect = f64.Rectangle{
			f64.Vec2{-math.MaxFloat32, -math.MaxFloat32},
			f64.Vec2{+math.MaxFloat32, +math.MaxFloat32},
		}
		window.LastFrameActive = current_frame
		window.IDStack = window.IDStack[:1]

		// Lock window rounding, border size and rounding so that altering the border sizes for children doesn't have side-effects.
		if flags&WindowFlagsChildWindow != 0 {
			window.WindowRounding = style.ChildRounding
		} else if flags&WindowFlagsPopup != 0 && flags&WindowFlagsModal == 0 {
			window.WindowRounding = style.PopupRounding
		} else {
			window.WindowRounding = style.WindowRounding
		}

		if flags&WindowFlagsChildWindow != 0 {
			window.WindowBorderSize = style.ChildBorderSize
		} else if flags&WindowFlagsPopup != 0 && flags&WindowFlagsModal == 0 {
			window.WindowBorderSize = style.PopupBorderSize
		} else {
			window.WindowBorderSize = style.WindowBorderSize
		}
		window.WindowPadding = style.WindowPadding

		if flags&WindowFlagsChildWindow != 0 && flags&(WindowFlagsAlwaysUseWindowPadding|WindowFlagsPopup) == 0 && window.WindowBorderSize == 0 {
			window.WindowPadding = f64.Vec2{0, 0}
			if flags&WindowFlagsMenuBar != 0 {
				window.WindowPadding.Y = style.WindowPadding.Y
			}
		}

		// Collapse window by double-clicking on title bar
		// At this point we don't have a clipping rectangle setup yet, so we can use the title bar area for hit detection and drawing
		if flags&WindowFlagsNoTitleBar == 0 && flags&WindowFlagsNoCollapse == 0 {
			title_bar_rect := window.TitleBarRect()
			if window.CollapseToggleWanted || (c.HoveredWindow == window && c.IsMouseHoveringRect(title_bar_rect.Min, title_bar_rect.Max) && c.IO.MouseDoubleClicked[0]) {
				window.Collapsed = !window.Collapsed
				c.MarkIniSettingsDirtyEx(window)
				c.FocusWindow(window)
			}
		} else {
			window.Collapsed = false
		}
		window.CollapseToggleWanted = false

		// SIZE

		// Update contents size from last frame for auto-fitting (unless explicitly specified)
		window.SizeContents = c.CalcSizeContents(window)

		// Hide popup/tooltip window when re-opening while we measure size (because we recycle the windows)
		if window.HiddenFrames > 0 {
			window.HiddenFrames--
		}

		if flags&(WindowFlagsPopup|WindowFlagsTooltip) != 0 && window_just_activated_by_user {
			window.HiddenFrames = 1
			if flags&WindowFlagsAlwaysAutoResize != 0 {
				if !window_size_x_set_by_api {
					window.Size.X = 0
					window.SizeFull.X = 0
				}
				if !window_size_y_set_by_api {
					window.Size.Y = 0
					window.SizeFull.Y = 0
				}
				window.SizeContents = f64.Vec2{0, 0}
			}
		}

		// Calculate auto-fit size, handle automatic resize
		size_auto_fit := c.CalcSizeAutoFit(window, window.SizeContents)
		size_full_modified := f64.Vec2{math.MaxFloat32, math.MaxFloat32}
		if flags&WindowFlagsAlwaysAutoResize != 0 && !window.Collapsed {
			// Using SetNextWindowSize() overrides ImGuiWindowFlags_AlwaysAutoResize, so it can be used on tooltips/popups, etc.
			if !window_size_x_set_by_api {
				size_full_modified.X = size_auto_fit.X
				window.SizeFull.X = size_full_modified.X
			}
			if !window_size_y_set_by_api {
				size_full_modified.Y = size_auto_fit.Y
				window.SizeFull.Y = size_full_modified.Y
			}
		} else if window.AutoFitFramesX > 0 || window.AutoFitFramesY > 0 {
			// Auto-fit only grows during the first few frames
			// We still process initial auto-fit on collapsed windows to get a window width, but otherwise don't honor ImGuiWindowFlags_AlwaysAutoResize when collapsed.
			if !window_size_x_set_by_api && window.AutoFitFramesX > 0 {
				if window.AutoFitOnlyGrows {
					size_full_modified.X = math.Max(window.SizeFull.X, size_auto_fit.X)
				} else {
					size_full_modified.X = size_auto_fit.X
				}
				window.SizeFull.X = size_full_modified.X
			}

			if !window_size_y_set_by_api && window.AutoFitFramesY > 0 {
				if window.AutoFitOnlyGrows {
					size_full_modified.Y = math.Max(window.SizeFull.Y, size_auto_fit.Y)
				} else {
					size_full_modified.Y = size_auto_fit.Y
				}
				window.SizeFull.Y = size_full_modified.Y
			}

			if !window.Collapsed {
				c.MarkIniSettingsDirtyEx(window)
			}
		}

		// Apply minimum/maximum window size constraints and final size
		window.SizeFull = c.CalcSizeAfterConstraint(window, window.SizeFull)
		window.Size = window.SizeFull
		if window.Collapsed && flags&WindowFlagsChildWindow == 0 {
			window.Size = window.TitleBarRect().Size()
		}

		// SCROLLBAR STATUS

		// Update scrollbar status (based on the Size that was effective during last frame or the auto-resized Size).
		if !window.Collapsed {
			// When reading the current size we need to read it after size constraints have been applied
			size_x_for_scrollbars := window.SizeFullAtLastBegin.X
			if size_full_modified.X != math.MaxFloat32 {
				size_x_for_scrollbars = window.SizeFull.X
			}

			size_y_for_scrollbars := window.SizeFullAtLastBegin.Y
			if size_full_modified.Y != math.MaxFloat32 {
				size_y_for_scrollbars = window.SizeFull.Y
			}

			window.ScrollbarY = flags&WindowFlagsAlwaysVerticalScrollbar != 0 || window.SizeContents.Y > size_y_for_scrollbars && flags&WindowFlagsNoScrollbar == 0
			if window.ScrollbarX && !window.ScrollbarY {
				window.ScrollbarY = (window.SizeContents.Y > size_y_for_scrollbars-style.ScrollbarSize) && flags&WindowFlagsNoScrollbar == 0
			}

			scrollbarSizeX := 0.0
			if window.ScrollbarY {
				scrollbarSizeX = style.ScrollbarSize
			}
			window.ScrollbarX = flags&WindowFlagsAlwaysHorizontalScrollbar != 0 ||
				((window.SizeContents.X > size_x_for_scrollbars-scrollbarSizeX) &&
					flags&WindowFlagsNoScrollbar == 0 && flags&WindowFlagsHorizontalScrollbar != 0)

			window.ScrollbarSizes = f64.Vec2{0, 0}
			if window.ScrollbarY {
				window.ScrollbarSizes.X = style.ScrollbarSize
			}
			if window.ScrollbarX {
				window.ScrollbarSizes.Y = style.ScrollbarSize
			}
		}

		// POSITION
		// Popup latch its initial position, will position itself when it appears next frame
		if window_just_activated_by_user {
			window.AutoPosLastDirection = DirNone
			if flags&WindowFlagsPopup != 0 && !window_pos_set_by_api {
				window.PosFloat = c.CurrentPopupStack[len(c.CurrentPopupStack)-1].OpenPopupPos
				window.Pos = window.PosFloat
			}
		}

		// Position child window
		if flags&WindowFlagsChildWindow != 0 {
			window.BeginOrderWithinParent = len(parent_window.DC.ChildWindows)
			parent_window.DC.ChildWindows = append(parent_window.DC.ChildWindows, window)
			if flags&WindowFlagsPopup == 0 && !window_pos_set_by_api && !window_is_child_tooltip {
				window.PosFloat = parent_window.DC.CursorPos
				window.Pos = window.PosFloat
			}
		}

		window_pos_with_pivot := (window.SetWindowPosVal.X != math.MaxFloat32 && window.HiddenFrames == 0)
		if window_pos_with_pivot {
			// Position given a pivot (e.g. for centering)
			windowPos := window.SizeFull.Scale2(window.SetWindowPosPivot)
			windowPos = window.SetWindowPosVal.Sub(windowPos)
			windowPos = windowPos.Max(style.DisplaySafeAreaPadding)
			c.SetWindowPos(window, windowPos, 0)
		} else if flags&WindowFlagsChildMenu != 0 {
			// Child menus typically request _any_ position within the parent menu item, and then our FindBestPopupWindowPos() function will move the new menu outside the parent bounds.
			// This is how we end up with child menus appearing (most-commonly) on the right of the parent menu.
			// We want some overlap to convey the relative depth of each popup (currently the amount of overlap it is hard-coded to style.ItemSpacing.x, may need to introduce another style value).
			horizontal_overlap := style.ItemSpacing.X
			parent_menu := parent_window_in_stack
			var rect_to_avoid f64.Rectangle
			if parent_menu.DC.MenuBarAppending {
				rect_to_avoid = f64.Rectangle{
					f64.Vec2{-math.MaxFloat32, parent_menu.Pos.Y + parent_menu.TitleBarHeight()},
					f64.Vec2{math.MaxFloat32, parent_menu.Pos.Y + parent_menu.TitleBarHeight() + parent_menu.MenuBarHeight()},
				}
			} else {
				rect_to_avoid = f64.Rectangle{
					f64.Vec2{parent_menu.Pos.X + horizontal_overlap, -math.MaxFloat32},
					f64.Vec2{parent_menu.Pos.X + parent_menu.Size.X - horizontal_overlap - parent_menu.ScrollbarSizes.X, math.MaxFloat32},
				}
			}
			window.PosFloat = c.FindBestWindowPosForPopup(window.PosFloat, window.Size, &window.AutoPosLastDirection, rect_to_avoid)
		} else if flags&WindowFlagsPopup != 0 && !window_pos_set_by_api && window_just_appearing_after_hidden_for_resize {
			rect_to_avoid := f64.Rectangle{
				f64.Vec2{window.PosFloat.X - 1, window.PosFloat.Y - 1},
				f64.Vec2{window.PosFloat.X + 1, window.PosFloat.Y + 1},
			}
			window.PosFloat = c.FindBestWindowPosForPopup(window.PosFloat, window.Size, &window.AutoPosLastDirection, rect_to_avoid)
		}

		// Position tooltip (always follows mouse)
		if flags&WindowFlagsTooltip != 0 && !window_pos_set_by_api && !window_is_child_tooltip {
			sc := c.Style.MouseCursorScale
			ref_pos := c.IO.MousePos
			if !c.NavDisableHighlight && c.NavDisableMouseHover {
				ref_pos = c.NavCalcPreferredMousePos()
			}
			var rect_to_avoid f64.Rectangle
			if !c.NavDisableHighlight && c.NavDisableMouseHover && c.IO.ConfigFlags&ConfigFlagsNavMoveMouse == 0 {
				rect_to_avoid = f64.Rectangle{f64.Vec2{ref_pos.X - 16, ref_pos.Y - 8}, f64.Vec2{ref_pos.X + 16, ref_pos.Y + 8}}
			} else {
				// FIXME: Hard-coded based on mouse cursor shape expectation. Exact dimension not very important.
				rect_to_avoid = f64.Rectangle{f64.Vec2{ref_pos.X - 16, ref_pos.Y - 8}, f64.Vec2{ref_pos.X + 24*sc, ref_pos.Y + 24*sc}}
			}
			window.PosFloat = c.FindBestWindowPosForPopup(ref_pos, window.Size, &window.AutoPosLastDirection, rect_to_avoid)
			if window.AutoPosLastDirection == DirNone {
				// If there's not enough room, for tooltip we prefer avoiding the cursor at all cost even if it means that part of the tooltip won't be visible.
				window.PosFloat = ref_pos.Add(f64.Vec2{2, 2})
			}
		}

		// Clamp position so it stays visible
		if flags&WindowFlagsChildWindow == 0 && flags&WindowFlagsTooltip == 0 {
			// Ignore zero-sized display explicitly to avoid losing positions if a window manager reports zero-sized window when initializing or minimizing.
			if !window_pos_set_by_api && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 && c.IO.DisplaySize.X > 0 && c.IO.DisplaySize.Y > 0 {
				padding := style.DisplayWindowPadding.Max(style.DisplaySafeAreaPadding)
				posFloat := window.PosFloat.Add(window.Size)
				posFloat = posFloat.Max(padding)
				posFloat = posFloat.Sub(window.Size)
				posFloat = posFloat.Min(c.IO.DisplaySize.Sub(padding))
				window.PosFloat = posFloat
			}
		}
		window.Pos = window.PosFloat.Floor()

		// Default item width. Make it proportional to window size if window manually resizes
		if window.Size.X > 0 && flags&WindowFlagsTooltip == 0 && flags&WindowFlagsAlwaysAutoResize == 0 {
			window.ItemWidthDefault = float64(int(window.Size.X * 0.65))
		} else {
			window.ItemWidthDefault = float64(int(c.FontSize * 16))
		}

		// Prepare for focus requests
		if window.FocusIdxAllRequestNext == math.MaxInt32 || window.FocusIdxAllCounter == -1 {
			window.FocusIdxAllRequestCurrent = math.MaxInt32
		} else {
			window.FocusIdxAllRequestCurrent = (window.FocusIdxAllRequestNext + (window.FocusIdxAllCounter + 1)) % (window.FocusIdxAllCounter + 1)
		}
		if window.FocusIdxTabRequestNext == math.MaxInt32 || window.FocusIdxTabCounter == -1 {
			window.FocusIdxTabRequestCurrent = math.MaxInt32
		} else {
			window.FocusIdxTabRequestCurrent = (window.FocusIdxTabRequestNext + (window.FocusIdxTabCounter + 1)) % (window.FocusIdxTabCounter + 1)
		}
		window.FocusIdxAllCounter = -1
		window.FocusIdxTabCounter = -1
		window.FocusIdxAllRequestNext = math.MaxInt32
		window.FocusIdxTabRequestNext = math.MaxInt32

		// Apply scrolling
		window.Scroll = c.CalcNextScrollFromScrollTargetAndClamp(window)
		window.ScrollTarget = f64.Vec2{math.MaxInt32, math.MaxInt32}

		// Apply focus, new windows appears in front
		want_focus := false
		if window_just_activated_by_user && flags&WindowFlagsNoFocusOnAppearing == 0 {
			if flags&(WindowFlagsChildWindow|WindowFlagsTooltip) == 0 || flags&WindowFlagsPopup != 0 {
				want_focus = true
			}
		}

		// Handle manual resize: Resize Grips, Borders, Gamepad
		border_held := -1
		resize_grip_col := [4]color.RGBA{}
		resize_grip_count := 1
		if flags&WindowFlagsResizeFromAnySide != 0 {
			resize_grip_count = 2
		}
		grip_draw_size := float64(int(math.Max(c.FontSize*1.35, window.WindowRounding+1.0+c.FontSize*0.2)))
		if !window.Collapsed {
			c.UpdateManualResize(window, size_auto_fit, &border_held, resize_grip_col[:resize_grip_count])
		}

		// DRAWING

		// Setup draw list and outer clipping rectangle
		window.DrawList.Clear()
		if c.Style.AntiAliasedLines {
			window.DrawList.Flags |= DrawListFlagsAntiAliasedLines
		}
		if c.Style.AntiAliasedFill {
			window.DrawList.Flags |= DrawListFlagsAntiAliasedFill
		}
		window.DrawList.PushTextureID(c.Font.ContainerAtlas.TexID)
		viewport_rect := c.GetViewportRect()
		if flags&WindowFlagsChildWindow != 0 && flags&WindowFlagsPopup == 0 && !window_is_child_tooltip {
			c.PushClipRect(parent_window.ClipRect.Min, parent_window.ClipRect.Max, true)
		} else {
			c.PushClipRect(viewport_rect.Min, viewport_rect.Max, true)
		}

		// Draw modal window background (darkens what is behind them)
		if flags&WindowFlagsModal != 0 && window == c.GetFrontMostModalRootWindow() {
			window.DrawList.AddRectFilled(
				viewport_rect.Min,
				viewport_rect.Max,
				c.GetColorFromStyleWithAlpha(ColModalWindowDarkening, c.ModalWindowDarkeningRatio),
			)
		}

		// Draw navigation selection/windowing rectangle background
		if c.NavWindowingTarget == window {
			bb := window.Rect()
			bb = bb.Expand(c.FontSize, c.FontSize)
			// Avoid drawing if the window covers all the viewport anyway
			if !viewport_rect.In(bb) {
				window.DrawList.AddRectFilledEx(
					bb.Min, bb.Max,
					c.GetColorFromStyleWithAlpha(ColNavWindowingHighlight, c.NavWindowingHighlightAlpha*0.25),
					c.Style.WindowRounding,
					DrawCornerFlagsAll,
				)
			}
		}

		// Draw window + handle manual resize
		window_rounding := window.WindowRounding
		window_border_size := window.WindowBorderSize
		title_bar_is_highlight := want_focus || (c.NavWindow != nil && window.RootWindowForTitleBarHighlight == c.NavWindow.RootWindowForTitleBarHighlight)
		title_bar_rect := window.TitleBarRect()
		if window.Collapsed {
			// Title bar only
			backup_border_size := style.FrameBorderSize
			c.Style.FrameBorderSize = window.WindowBorderSize

			var title_bar_col color.RGBA
			if title_bar_is_highlight && !c.NavDisableHighlight {
				title_bar_col = c.GetColorFromStyle(ColTitleBgActive)
			} else {
				title_bar_col = c.GetColorFromStyle(ColTitleBgCollapsed)
			}

			c.RenderFrameEx(title_bar_rect.Min, title_bar_rect.Max, title_bar_col, true, window_rounding)
			c.Style.FrameBorderSize = backup_border_size
		} else {
			// Window background
			bg_col := c.GetColorFromStyle(c.GetWindowBgColorIdxFromFlags(flags))
			if c.NextWindowData.BgAlphaCond != 0 {
				bg_col.A = uint8(c.NextWindowData.BgAlphaVal * 255)
				c.NextWindowData.BgAlphaCond = 0
			}
			drawCornerFlags := DrawCornerFlagsBot
			if flags&WindowFlagsNoTitleBar != 0 {
				drawCornerFlags = DrawCornerFlagsAll
			}
			window.DrawList.AddRectFilledEx(
				window.Pos.Add(f64.Vec2{0, window.TitleBarHeight()}),
				window.Pos.Add(window.Size),
				bg_col, window_rounding, drawCornerFlags,
			)

			// Title bar
			var title_bar_col color.RGBA
			if window.Collapsed {
				title_bar_col = c.GetColorFromStyle(ColTitleBgCollapsed)
			} else if title_bar_is_highlight {
				title_bar_col = c.GetColorFromStyle(ColTitleBgActive)
			} else {
				title_bar_col = c.GetColorFromStyle(ColTitleBg)
			}
			if flags&WindowFlagsNoTitleBar == 0 {
				window.DrawList.AddRectFilledEx(title_bar_rect.Min, title_bar_rect.Max, title_bar_col, window_rounding, DrawCornerFlagsTop)
			}

			// Menu bar
			if flags&WindowFlagsMenuBar != 0 {
				menu_bar_rect := window.MenuBarRect()
				menu_bar_rect = menu_bar_rect.Intersect(window.Rect())

				rounding := 0.0
				if flags&WindowFlagsNoTitleBar != 0 {
					rounding = window_rounding
				}
				window.DrawList.AddRectFilledEx(
					menu_bar_rect.Min, menu_bar_rect.Max,
					c.GetColorFromStyle(ColMenuBarBg), rounding, DrawCornerFlagsTop,
				)

				if style.FrameBorderSize > 0 && menu_bar_rect.Max.Y < window.Pos.Y+window.Size.Y {
					window.DrawList.AddLineEx(
						menu_bar_rect.BL(), menu_bar_rect.BR(),
						c.GetColorFromStyle(ColBorder), style.FrameBorderSize,
					)
				}
			}

			// Scrollbars
			if window.ScrollbarX {
				c.Scrollbar(LayoutTypeHorizontal)
			}
			if window.ScrollbarY {
				c.Scrollbar(LayoutTypeVertical)
			}

			// Render resize grips (after their input handling so we don't have a frame of latency)
			if flags&WindowFlagsNoResize == 0 {
				for resize_grip_n := range resize_grip_def {
					grip := resize_grip_def[resize_grip_n]
					corner := window.Pos.Lerp2(grip.CornerPos, window.Pos.Add(window.Size))

					var l1, l2 f64.Vec2
					if resize_grip_n&1 != 0 {
						l1 = f64.Vec2{window_border_size, grip_draw_size}
						l2 = f64.Vec2{grip_draw_size, window_border_size}
					} else {
						l1 = f64.Vec2{grip_draw_size, window_border_size}
						l2 = f64.Vec2{window_border_size, grip_draw_size}
					}
					p1 := grip.InnerDir.Scale2(l1)
					p2 := grip.InnerDir.Scale2(l2)
					p1 = corner.Add(p1)
					p2 = corner.Add(p2)
					window.DrawList.PathLineTo(p1)
					window.DrawList.PathLineTo(p2)
					window.DrawList.PathArcToFast(
						f64.Vec2{
							corner.X + grip.InnerDir.X*(window_rounding+window_border_size),
							corner.Y + grip.InnerDir.Y*(window_rounding+window_border_size),
						},
						window_rounding, grip.AngleMin12, grip.AngleMax12,
					)
					window.DrawList.PathFillConvex(resize_grip_col[resize_grip_n])
				}
			}

			// Borders
			if window_border_size > 0 {
				window.DrawList.AddRectEx(window.Pos, window.Pos.Add(window.Size), c.GetColorFromStyle(ColBorder), window_rounding, DrawCornerFlagsAll, window_border_size)
			}
			if border_held != -1 {
				border := c.GetBorderRect(window, border_held, grip_draw_size, 0.0)
				window.DrawList.AddLineEx(border.Min, border.Max, c.GetColorFromStyle(ColSeparatorActive), math.Max(1.0, window_border_size))
			}
			if style.FrameBorderSize > 0 && flags&WindowFlagsNoTitleBar == 0 {
				window.DrawList.AddLineEx(
					title_bar_rect.BL().Add(f64.Vec2{style.WindowBorderSize, -1}),
					title_bar_rect.BR().Add(f64.Vec2{-style.WindowBorderSize, -1}),
					c.GetColorFromStyle(ColBorder), style.FrameBorderSize,
				)
			}
		}

		// Draw navigation selection/windowing rectangle border
		if c.NavWindowingTarget == window {
			rounding := math.Max(window.WindowRounding, c.Style.WindowRounding)
			bb := window.Rect()
			bb = bb.Expand(c.FontSize, c.FontSize)
			// If a window fits the entire viewport, adjust its highlight inward
			if viewport_rect.In(bb) {
				bb = bb.Expand(-c.FontSize-1.0, -c.FontSize-1.0)
				rounding = window.WindowRounding
			}
			window.DrawList.AddRectEx(bb.Min, bb.Max, c.GetColorFromStyleWithAlpha(ColNavWindowingHighlight, c.NavWindowingHighlightAlpha), rounding, ^0, 3.0)
		}

		// Store a backup of SizeFull which we will use next frame to decide if we need scrollbars.
		window.SizeFullAtLastBegin = window.SizeFull

		// Update ContentsRegionMax. All the variable it depends on are set above in this function.
		window.ContentsRegionRect.Min.X = -window.Scroll.X + window.WindowPadding.X
		window.ContentsRegionRect.Min.Y = -window.Scroll.Y + window.WindowPadding.Y + window.TitleBarHeight() + window.MenuBarHeight()
		if window.SizeContentsExplicit.X != 0.0 {
			window.ContentsRegionRect.Max.X = -window.Scroll.X - window.WindowPadding.X + window.SizeContentsExplicit.X
		} else {
			window.ContentsRegionRect.Max.X = -window.Scroll.X - window.WindowPadding.X + window.Size.X - window.ScrollbarSizes.X
		}
		if window.SizeContentsExplicit.Y != 0.0 {
			window.ContentsRegionRect.Max.Y = -window.Scroll.Y - window.WindowPadding.Y + window.SizeContentsExplicit.Y
		} else {
			window.ContentsRegionRect.Max.Y = -window.Scroll.Y - window.WindowPadding.Y + window.Size.Y - window.ScrollbarSizes.Y
		}

		// Setup drawing context
		// (NB: That term "drawing context / DC" lost its meaning a long time ago. Initially was meant to hold transient data only. Nowadays difference between window-> and window->DC-> is dubious.)
		window.DC.IndentX = 0.0 + window.WindowPadding.X - window.Scroll.X
		window.DC.GroupOffsetX = 0.0
		window.DC.ColumnsOffsetX = 0.0
		window.DC.CursorStartPos = window.Pos.Add(f64.Vec2{
			window.DC.IndentX + window.DC.ColumnsOffsetX,
			window.TitleBarHeight() + window.MenuBarHeight() + window.WindowPadding.Y - window.Scroll.Y,
		})
		window.DC.CursorPos = window.DC.CursorStartPos
		window.DC.CursorPosPrevLine = window.DC.CursorPos
		window.DC.CursorMaxPos = window.DC.CursorStartPos
		window.DC.CurrentLineHeight = 0
		window.DC.PrevLineHeight = 0
		window.DC.CurrentLineTextBaseOffset = 0
		window.DC.PrevLineTextBaseOffset = 0
		window.DC.NavHideHighlightOneFrame = false
		window.DC.NavHasScroll = c.GetScrollMaxY() > 0
		window.DC.NavLayerActiveMask = window.DC.NavLayerActiveMaskNext
		window.DC.NavLayerActiveMaskNext = 0x00
		window.DC.MenuBarAppending = false
		window.DC.MenuBarOffsetX = math.Max(window.WindowPadding.X, style.ItemSpacing.X)
		window.DC.LogLinePosY = window.DC.CursorPos.Y - 9999
		window.DC.ChildWindows = window.DC.ChildWindows[:0]
		window.DC.LayoutType = LayoutTypeVertical
		window.DC.ParentLayoutType = LayoutTypeVertical
		if parent_window != nil {
			window.DC.ParentLayoutType = parent_window.DC.LayoutType
		}
		window.DC.ItemFlags = ItemFlagsDefault_
		window.DC.ItemWidth = window.ItemWidthDefault
		window.DC.TextWrapPos = -1.0 // disabled
		window.DC.ItemFlagsStack = window.DC.ItemFlagsStack[:0]
		window.DC.ItemWidthStack = window.DC.ItemWidthStack[:0]
		window.DC.TextWrapPosStack = window.DC.TextWrapPosStack[:0]
		window.DC.ColumnsSet = nil
		window.DC.TreeDepth = 0
		window.DC.TreeDepthMayJumpToParentOnPop = 0x00
		window.DC.StateStorage = window.StateStorage
		window.DC.GroupStack = window.DC.GroupStack[:0]
		window.MenuColumns.Update(3, style.ItemSpacing.X, window_just_activated_by_user)

		if flags&WindowFlagsChildWindow != 0 && window.DC.ItemFlags != parent_window.DC.ItemFlags {
			window.DC.ItemFlags = parent_window.DC.ItemFlags
			window.DC.ItemFlagsStack = append(window.DC.ItemFlagsStack, window.DC.ItemFlags)
		}

		if window.AutoFitFramesX > 0 {
			window.AutoFitFramesX--
		}
		if window.AutoFitFramesY > 0 {
			window.AutoFitFramesY--
		}

		// Apply focus (we need to call FocusWindow() AFTER setting DC.CursorStartPos so our initial navigation reference rectangle can start around there)
		if want_focus {
			c.FocusWindow(window)
			c.NavInitWindow(window, false)
		}

		// Title bar
		if flags&WindowFlagsNoTitleBar == 0 {
			// Close & collapse button are on layer 1 (same as menus) and don't default focus
			item_flags_backup := window.DC.ItemFlags
			window.DC.ItemFlags |= ItemFlagsNoNavDefaultFocus
			window.DC.NavLayerCurrent++
			window.DC.NavLayerCurrentMask <<= 1

			// Collapse button
			if flags&WindowFlagsNoCollapse == 0 {
				id := window.GetID("#COLLAPSE")
				one := f64.Vec2{1, 1}

				bx := style.FramePadding.Add(one)
				bx = window.Pos.Add(bx)
				by := f64.Vec2{c.FontSize, c.FontSize}
				by = by.Add(one)
				by = style.FramePadding.Add(by)
				by = window.Pos.Add(by)
				bb := f64.Rectangle{bx, by}

				// To allow navigation
				c.ItemAdd(bb, id)

				_, _, pressed := c.ButtonBehavior(bb, id, 0)
				if pressed {
					// Defer collapsing to next frame as we are too far in the Begin() function
					window.CollapseToggleWanted = true
				}
				c.RenderNavHighlight(bb, id)
				dir := DirDown
				if window.Collapsed {
					dir = DirRight
				}
				c.RenderArrowEx(window.Pos.Add(style.FramePadding), dir, 1)
			}

			// Close button
			if p_open != nil {
				pad := style.FramePadding.Y
				rad := c.FontSize * 0.5
				padRect := f64.Vec2{-pad - rad, pad + rad}
				if c.CloseButton(window.GetID("#CLOSE"), window.Rect().TR().Add(padRect), rad+1) {
					*p_open = false
				}
			}

			window.DC.NavLayerCurrent--
			window.DC.NavLayerCurrentMask >>= 1
			window.DC.ItemFlags = item_flags_backup

			// Title text (FIXME: refactor text alignment facilities along with RenderText helpers)
			text_size := c.CalcTextSizeEx(name, true, -1)
			text_r := title_bar_rect
			pad_left := style.FramePadding.X
			if flags&WindowFlagsNoCollapse == 0 {
				pad_left = style.FramePadding.X + c.FontSize + style.ItemInnerSpacing.X
			}
			pad_right := style.FramePadding.X
			if p_open != nil {
				pad_right = style.FramePadding.X + c.FontSize + style.ItemInnerSpacing.X
			}
			if style.WindowTitleAlign.X > 0 {
				pad_right = f64.Lerp(style.WindowTitleAlign.X, pad_right, pad_left)
			}
			text_r.Min.X += pad_left
			text_r.Max.X -= pad_right
			clip_rect := text_r
			clip_rect.Max.X = window.Pos.X + window.Size.X
			if p_open != nil {
				clip_rect.Max.X -= title_bar_rect.Dy() - 3
			} else {
				clip_rect.Max.X -= style.FramePadding.X
			}
			c.RenderTextClippedEx(text_r.Min, text_r.Max, name, &text_size, style.WindowTitleAlign, &clip_rect)
		}

		// Save clipped aabb so we can access it in constant-time in FindHoveredWindow()
		window.WindowRectClipped = window.Rect()
		window.WindowRectClipped = window.WindowRectClipped.Intersect(window.ClipRect)

		// Pressing CTRL+C while holding on a window copy its content to the clipboard
		// This works but 1. doesn't handle multiple Begin/End pairs, 2. recursing into another Begin/End pair - so we need to work that out and add better logging scope.
		// Maybe we can support CTRL+C on every element?

		// Inner rectangle
		// We set this up after processing the resize grip so that our clip rectangle doesn't lag by a frame
		// Note that if our window is collapsed we will end up with a null clipping rectangle which is the correct behavior.
		window.InnerRect.Min.X = title_bar_rect.Min.X + window.WindowBorderSize
		window.InnerRect.Min.Y = title_bar_rect.Max.Y + window.MenuBarHeight()
		if flags&WindowFlagsMenuBar != 0 || flags&WindowFlagsNoTitleBar == 0 {
			window.InnerRect.Min.Y += style.FrameBorderSize
		} else {
			window.InnerRect.Min.Y += window.WindowBorderSize
		}
		window.InnerRect.Max.X = window.Pos.X + window.Size.X - window.ScrollbarSizes.X - window.WindowBorderSize
		window.InnerRect.Max.Y = window.Pos.Y + window.Size.Y - window.ScrollbarSizes.Y - window.WindowBorderSize

		// Inner clipping rectangle
		// Force round operator last to ensure that e.g. (int)(max.x-min.x) in user's render code produce correct result.
		window.InnerClipRect.Min.X = math.Floor(0.5 + window.InnerRect.Min.X + math.Max(0, math.Floor(window.WindowPadding.X*0.5-window.WindowBorderSize)))
		window.InnerClipRect.Min.Y = math.Floor(0.5 + window.InnerRect.Min.Y)
		window.InnerClipRect.Max.X = math.Floor(0.5 + window.InnerRect.Max.X - math.Max(0, math.Floor(window.WindowPadding.X*0.5-window.WindowBorderSize)))
		window.InnerClipRect.Max.Y = math.Floor(0.5 + window.InnerRect.Max.Y)

		// After Begin() we fill the last item / hovered data using the title bar data. Make that a standard behavior (to allow usage of context menus on title bar only, etc.).
		window.DC.LastItemId = window.MoveId
		window.DC.LastItemStatusFlags = 0
		if c.IsMouseHoveringRectEx(title_bar_rect.Min, title_bar_rect.Max, false) {
			window.DC.LastItemStatusFlags = ItemStatusFlagsHoveredRect
		}
		window.DC.LastItemRect = title_bar_rect
	}

	c.PushClipRect(window.InnerClipRect.Min, window.InnerClipRect.Max, true)

	// Clear 'accessed' flag last thing (After PushClipRect which will set the flag. We want the flag to stay false when the default "Debug" window is unused)
	if first_begin_of_the_frame {
		window.WriteAccessed = false
	}

	window.BeginCount++
	c.NextWindowData.SizeConstraintCond = 0

	// Child window can be out of sight and have "negative" clip windows.
	// Mark them as collapsed so commands are skipped earlier (we can't manually collapse because they have no title bar).
	if flags&WindowFlagsChildWindow != 0 {
		window.Collapsed = parent_window != nil && parent_window.Collapsed

		if flags&WindowFlagsAlwaysAutoResize == 0 && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0 {
			if window.WindowRectClipped.Min.X >= window.WindowRectClipped.Max.X ||
				window.WindowRectClipped.Min.Y >= window.WindowRectClipped.Max.Y {
				window.Collapsed = true
			}
		}

		// We also hide the window from rendering because we've already added its border to the command list.
		// (we could perform the check earlier in the function but it is simpler at this point)
		if window.Collapsed {
			window.Active = false
		}
	}
	if style.Alpha <= 0.0 {
		window.Active = false
	}

	// Return false if we don't intend to display anything to allow user to perform an early out optimization
	window.SkipItems = (window.Collapsed || !window.Active) && window.AutoFitFramesX <= 0 && window.AutoFitFramesY <= 0
	return !window.SkipItems
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
	var (
		offset, size f64.Vec2
		uv           [4]f64.Vec2
	)
	if c.IO.MouseDrawCursor && c.IO.Fonts.GetMouseCursorTexData(c.MouseCursor, &offset, &size, uv[:], uv[:]) {
		pos := c.IO.MousePos.Sub(offset)
		tex_id := c.IO.Fonts.TexID
		sc := c.Style.MouseCursorScale
		c.OverlayDrawList.PushTextureID(tex_id)

		// Shadow
		x := f64.Vec2{1, 0}
		x = x.Scale(sc)
		x = pos.Add(x)
		y := f64.Vec2{1, 0}
		y = y.Scale(sc)
		y = pos.Add(y)
		s := size.Scale(sc)
		y = y.Add(s)
		c.OverlayDrawList.AddImageEx(tex_id, x, y, uv[2], uv[3], color.RGBA{0, 0, 0, 48})

		// Shadow
		x = f64.Vec2{2, 0}
		x = x.Scale(sc)
		x = pos.Add(x)
		y = f64.Vec2{2, 0}
		y = y.Scale(sc)
		y = pos.Add(y)
		s = size.Scale(sc)
		y = y.Add(s)
		c.OverlayDrawList.AddImageEx(tex_id, x, y, uv[2], uv[3], color.RGBA{0, 0, 0, 48})

		// Black border
		x = pos
		y = pos.Add(s)
		c.OverlayDrawList.AddImageEx(tex_id, x, y, uv[2], uv[3], color.RGBA{0, 0, 0, 255})

		// White fill
		x = pos
		y = pos.Add(s)
		c.OverlayDrawList.AddImageEx(tex_id, x, y, uv[0], uv[1], color.RGBA{255, 255, 255, 255})
	}

	if len(c.OverlayDrawList.VtxBuffer) == 0 {
		c.AddDrawListToDrawData(&c.DrawDataBuilder.Layers[0], &c.OverlayDrawList)
	}

	// Setup ImDrawData structure for end-user
	c.SetupDrawData(c.DrawDataBuilder.Layers[0], &c.DrawData)
	c.IO.MetricsRenderVertices = c.DrawData.TotalVtxCount
	c.IO.MetricsRenderIndices = c.DrawData.TotalIdxCount
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
	c.RenderNavHighlightEx(bb, id, NavHighlightFlagsTypeDefault)
}

func (c *Context) RenderNavHighlightEx(bb f64.Rectangle, id ID, flags NavHighlightFlags) {
	if id != c.NavId {
		return
	}
	if c.NavDisableHighlight && flags&NavHighlightFlagsAlwaysDraw == 0 {
		return
	}
	window := c.GetCurrentWindow()
	if window.DC.NavHideHighlightOneFrame {
		return
	}

	rounding := 0.0
	if flags&NavHighlightFlagsNoRounding != 0 {
		rounding = c.Style.FrameRounding
	}
	display_rect := bb
	display_rect = display_rect.Intersect(window.ClipRect)
	if flags&NavHighlightFlagsTypeDefault != 0 {
		const THICKNESS = 2.0
		const DISTANCE = 3.0 + THICKNESS*0.5
		display_rect = display_rect.Expand(DISTANCE, DISTANCE)
		fully_visible := display_rect.In(window.ClipRect)
		if !fully_visible {
			window.DrawList.PushClipRect(display_rect.Min, display_rect.Max)
		}
		window.DrawList.AddRectEx(
			display_rect.Min.Add(f64.Vec2{THICKNESS * 0.5, THICKNESS * 0.5}),
			display_rect.Max.Sub(f64.Vec2{THICKNESS * 0.5, THICKNESS * 0.5}),
			c.GetColorFromStyle(ColNavHighlight),
			rounding, DrawCornerFlagsAll, THICKNESS,
		)
		if !fully_visible {
			window.DrawList.PopClipRect()
		}
	}
	if flags&NavHighlightFlagsTypeThin != 0 {
		window.DrawList.AddRectEx(display_rect.Min, display_rect.Max, c.GetColorFromStyle(ColNavHighlight), rounding, ^0, 1.0)
	}
}

func (c *Context) RenderBullet(pos f64.Vec2) {
	window := c.CurrentWindow
	window.DrawList.AddCircleFilledEx(pos, c.FontSize*0.20, c.GetColorFromStyle(ColText), 8)
}

func (c *Context) RenderCheckMark(pos f64.Vec2, col color.RGBA, sz float64) {
	window := c.CurrentWindow

	thickness := math.Max(sz/5.0, 1.0)
	sz -= thickness * 0.5
	pos = pos.Add(f64.Vec2{thickness * 0.25, thickness * 0.25})

	third := sz / 3.0
	bx := pos.X + third
	by := pos.Y + sz - third*0.5
	window.DrawList.PathLineTo(f64.Vec2{bx - third, by - third})
	window.DrawList.PathLineTo(f64.Vec2{bx, by})
	window.DrawList.PathLineTo(f64.Vec2{bx + third*2, by - third*2})
	window.DrawList.PathStrokeEx(col, false, thickness)
}

func (c *Context) RenderFrame(p_min, p_max f64.Vec2, col color.RGBA) {
	c.RenderFrameEx(p_min, p_max, col, true, 0)
}

func (c *Context) RenderFrameEx(p_min, p_max f64.Vec2, col color.RGBA, border bool, rounding float64) {
}

func (c *Context) RenderArrow(pos f64.Vec2, dir Dir) {
	c.RenderArrowEx(pos, dir, 1)
}

func (c *Context) RenderArrowEx(pos f64.Vec2, dir Dir, scale float64) {
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

func (d *DrawList) AddRect(p_min, p_max f64.Vec2, col color.RGBA) {
	d.AddRectEx(p_min, p_max, col, 0, DrawCornerFlagsAll, 1)
}

func (d *DrawList) AddRectEx(p_min, p_max f64.Vec2, col color.RGBA, rounding float64, rounding_corner_flags DrawCornerFlags, thickness float64) {
}

func (d *DrawList) AddRectFilled(p_min, p_max f64.Vec2, col color.RGBA) {
	d.AddRectFilledEx(p_min, p_max, col, 0, DrawCornerFlagsAll)
}

func (d *DrawList) AddRectFilledEx(p_min, p_max f64.Vec2, col color.RGBA, rounding float64, rounding_corners_flags DrawCornerFlags) {
}

func (d *DrawList) AddCircleFilled(centre f64.Vec2, radius float64, col color.RGBA) {
	d.AddCircleFilledEx(centre, radius, col, 9)
}

func (d *DrawList) AddCircleFilledEx(centre f64.Vec2, radius float64, col color.RGBA, num_segments int) {
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

func (d *DrawList) PathArcToFast(centre f64.Vec2, radius float64, a_min_of_12, a_max_of_12 int) {
}

func (d *DrawList) PathFillConvex(col color.RGBA) {
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

// Handle resize for: Resize Grips, Borders, Gamepad
// TODO
func (c *Context) UpdateManualResize(window *Window, size_auto_fit f64.Vec2, border_held *int, resize_grip_col []color.RGBA) {
	flags := window.Flags
	if flags&WindowFlagsNoResize != 0 || flags&WindowFlagsAlwaysAutoResize != 0 || window.AutoFitFramesX > 0 || window.AutoFitFramesY > 0 {
		return
	}

	resize_border_count := 0
	if flags&WindowFlagsResizeFromAnySide != 0 {
		resize_border_count = 4
	}
	grip_draw_size := float64(int(math.Max(c.FontSize*1.35, window.WindowRounding+1.0+c.FontSize*0.2)))
	grip_hover_size := float64(int(grip_draw_size * 0.75))

	pos_target := f64.Vec2{math.MaxFloat32, math.MaxFloat32}
	size_target := f64.Vec2{math.MaxFloat32, math.MaxFloat32}

	_, _, _ = resize_border_count, pos_target, size_target

	// Manual resize grips
	c.PushID("#RESIZE")
	for resize_grip_n := range resize_grip_col {
		grip := resize_grip_def[resize_grip_n]
		corner := window.Pos.Lerp2(grip.CornerPos, window.Pos.Add(window.Size))

		// Using the FlattenChilds button flag we make the resize button accessible even if we are hovering over a child window
		resize_rect := f64.Rectangle{
			corner,
			corner.Add(grip.InnerDir.Scale(grip_hover_size)),
		}
		resize_rect = resize_rect.Canon()

		hovered, held, _ := c.ButtonBehavior(resize_rect, window.GetIDByInt(resize_grip_n), ButtonFlagsFlattenChildren|ButtonFlagsNoNavFocus)
		if hovered || held {
			if resize_grip_n&1 != 0 {
				c.MouseCursor = MouseCursorResizeNESW
			} else {
				c.MouseCursor = MouseCursorResizeNWSE
			}
		}

		if c.HoveredWindow == window && held && c.IO.MouseDoubleClicked[0] && resize_grip_n == 0 {
			// Manual auto-fit when double-clicking
			size_target = c.CalcSizeAfterConstraint(window, size_auto_fit)
			c.ClearActiveID()
		} else if held {
			// Resize from any of the four corners
			// We don't use an incremental MouseDelta but rather compute an absolute target size based on mouse position
		}

		if resize_grip_n == 0 || held || hovered {
			cornerTarget := resize_rect.Size()
			cornerTarget = cornerTarget.Scale2(grip.CornerPos)
		}
	}
}

func (c *Context) CalcNextScrollFromScrollTargetAndClamp(window *Window) f64.Vec2 {
	scroll := window.Scroll
	cr_x := window.ScrollTargetCenterRatio.X
	cr_y := window.ScrollTargetCenterRatio.Y
	if window.ScrollTarget.X < math.MaxFloat32 {
		scroll.X = window.ScrollTarget.X - cr_x*(window.SizeFull.X-window.ScrollbarSizes.X)
	}
	if window.ScrollTarget.Y < math.MaxFloat32 {
		scroll.Y = window.ScrollTarget.Y - (1-cr_y)*(window.TitleBarHeight()+window.MenuBarHeight()) - cr_y*(window.SizeFull.Y-window.ScrollbarSizes.Y)
	}
	scroll = scroll.Max(f64.Vec2{0, 0})
	if !window.Collapsed && !window.SkipItems {
		scroll.X = math.Min(scroll.X, c.GetWindowScrollMaxX(window))
		scroll.Y = math.Min(scroll.Y, c.GetWindowScrollMaxY(window))
	}

	return scroll
}

// Vertical scrollbar
// The entire piece of code below is rather confusing because:
// - We handle absolute seeking (when first clicking outside the grab) and relative manipulation (afterward or when clicking inside the grab)
// - We store values as normalized ratio and in a form that allows the window content to change while we are holding on a scrollbar
// - We handle both horizontal and vertical scrollbars, which makes the terminology not ideal.
func (c *Context) Scrollbar(direction LayoutType) {
	window := c.CurrentWindow

	horizontal := direction == LayoutTypeHorizontal
	style := &c.Style
	var id ID
	if horizontal {
		id = window.GetID("#SCROLLX")
	} else {
		id = window.GetID("#SCROLLY")
	}

	// Render background
	other_scrollbar := window.ScrollbarY
	if !horizontal {
		other_scrollbar = window.ScrollbarX
	}
	other_scrollbar_size_w := 0.0
	if other_scrollbar {
		other_scrollbar_size_w = style.ScrollbarSize
	}
	window_rect := window.Rect()
	border_size := window.WindowBorderSize

	var bb f64.Rectangle
	if horizontal {
		bb = f64.Rect(
			window.Pos.X+border_size,
			window_rect.Max.Y-style.ScrollbarSize,
			window_rect.Max.X-other_scrollbar_size_w-border_size,
			window_rect.Max.Y-border_size,
		)
	} else {
		bb = f64.Rect(
			window_rect.Max.X-style.ScrollbarSize,
			window.Pos.Y+border_size,
			window_rect.Max.X-border_size,
			window_rect.Max.Y-other_scrollbar_size_w-border_size,
		)
	}
	if !horizontal {
		bb.Min.Y += window.TitleBarHeight()
		if window.Flags&WindowFlagsMenuBar != 0 {
			bb.Min.Y += window.MenuBarHeight()
		}
	}
	if bb.Dx() <= 0 || bb.Dy() <= 0 {
		return
	}

	var window_rounding_corners DrawCornerFlags
	if horizontal {
		window_rounding_corners = DrawCornerFlagsBotLeft
		if !other_scrollbar {
			window_rounding_corners |= DrawCornerFlagsBotRight
		}
	} else {
		if window.Flags&WindowFlagsNoTitleBar != 0 && window.Flags&WindowFlagsMenuBar == 0 {
			window_rounding_corners |= DrawCornerFlagsTopRight
		}
		if !other_scrollbar {
			window_rounding_corners |= DrawCornerFlagsBotRight
		}
	}
	window.DrawList.AddRectFilledEx(bb.Min, bb.Max, c.GetColorFromStyle(ColScrollbarBg), window.WindowRounding, window_rounding_corners)
	bb = bb.Expand(
		-f64.Clamp(float64(int((bb.Max.X-bb.Min.X-2.0)*0.5)), 0, 3),
		-f64.Clamp(float64(int((bb.Max.Y-bb.Min.Y-2.0)*0.5)), 0, 3),
	)

	// V denote the main, longer axis of the scrollbar (= height for a vertical scrollbar)
	var scrollbar_size_v, scroll_v, win_size_avail_v, win_size_contents_v float64
	if horizontal {
		scrollbar_size_v = bb.Dx()
		scroll_v = window.Scroll.X
		win_size_avail_v = window.SizeFull.X - other_scrollbar_size_w
		win_size_contents_v = window.SizeContents.X
	} else {
		scrollbar_size_v = bb.Dy()
		scroll_v = window.Scroll.Y
		win_size_avail_v = window.SizeFull.Y - other_scrollbar_size_w
		win_size_contents_v = window.SizeContents.Y
	}

	// Calculate the height of our grabbable box. It generally represent the amount visible (vs the total scrollable amount)
	// But we maintain a minimum size in pixel to allow for the user to still aim inside.
	win_size_v := math.Max(math.Max(win_size_contents_v, win_size_avail_v), 1.0)
	grab_h_pixels := f64.Clamp(scrollbar_size_v*(win_size_avail_v/win_size_v), style.GrabMinSize, scrollbar_size_v)
	grab_h_norm := grab_h_pixels / scrollbar_size_v

	// Handle input right away. None of the code of Begin() is relying on scrolling position before calling Scrollbar().
	previously_held := c.ActiveId == id
	hovered, held, _ := c.ButtonBehavior(bb, id, ButtonFlagsNoNavFocus)

	scroll_max := math.Max(1.0, win_size_contents_v-win_size_avail_v)
	scroll_ratio := f64.Saturate(scroll_v / scroll_max)
	grab_v_norm := scroll_ratio * (scrollbar_size_v - grab_h_pixels) / scrollbar_size_v
	if held && grab_h_norm < 1.0 {
		var (
			scrollbar_pos_v, mouse_pos_v float64
			click_delta_to_grab_center_v *float64
		)
		if horizontal {
			scrollbar_pos_v = bb.Min.X
			mouse_pos_v = c.IO.MousePos.X
			click_delta_to_grab_center_v = &c.ScrollbarClickDeltaToGrabCenter.X
		} else {
			scrollbar_pos_v = bb.Min.Y
			mouse_pos_v = c.IO.MousePos.Y
			click_delta_to_grab_center_v = &c.ScrollbarClickDeltaToGrabCenter.Y
		}

		// Click position in scrollbar normalized space (0.0f->1.0f)
		clicked_v_norm := f64.Saturate((mouse_pos_v - scrollbar_pos_v) / scrollbar_size_v)
		c.SetHoveredID(id)

		seek_absolute := false
		if !previously_held {
			// On initial click calculate the distance between mouse and the center of the grab
			if clicked_v_norm >= grab_v_norm && clicked_v_norm <= grab_v_norm+grab_h_norm {
				*click_delta_to_grab_center_v = clicked_v_norm - grab_v_norm - grab_h_norm*0.5
			} else {
				seek_absolute = true
				*click_delta_to_grab_center_v = 0
			}
		}

		// Apply scroll
		// It is ok to modify Scroll here because we are being called in Begin() after the calculation of SizeContents and before setting up our starting position
		scroll_v_norm := f64.Saturate((clicked_v_norm - *click_delta_to_grab_center_v - grab_h_norm*0.5) / (1.0 - grab_h_norm))
		scroll_v = float64(int(0.5 + scroll_v_norm*scroll_max))
		if horizontal {
			window.Scroll.X = scroll_v
		} else {
			window.Scroll.Y = scroll_v
		}

		// Update values for rendering
		scroll_ratio = f64.Saturate(scroll_v / scroll_max)
		grab_v_norm = scroll_ratio * (scrollbar_size_v - grab_h_pixels) / scrollbar_size_v

		// Update distance to grab now that we have seeked and saturated
		if seek_absolute {
			*click_delta_to_grab_center_v = clicked_v_norm - grab_v_norm - grab_h_norm*0.5
		}
	}

	// Render
	var grab_rect f64.Rectangle
	var grab_col color.RGBA
	switch {
	case held:
		grab_col = c.GetColorFromStyle(ColScrollbarGrabActive)
	case hovered:
		grab_col = c.GetColorFromStyle(ColScrollbarGrabHovered)
	default:
		grab_col = c.GetColorFromStyle(ColScrollbarGrab)
	}

	if horizontal {
		grab_rect = f64.Rectangle{
			f64.Vec2{
				f64.Lerp(grab_v_norm, bb.Min.X, bb.Max.X),
				bb.Min.Y,
			},
			f64.Vec2{
				math.Min(f64.Lerp(grab_v_norm, bb.Min.X, bb.Max.X)+grab_h_pixels, window_rect.Max.X),
				bb.Max.Y,
			},
		}
	} else {
		grab_rect = f64.Rectangle{
			f64.Vec2{
				bb.Min.X,
				f64.Lerp(grab_v_norm, bb.Min.Y, bb.Max.Y),
			},
			f64.Vec2{
				bb.Max.X,
				math.Min(f64.Lerp(grab_v_norm, bb.Min.Y, bb.Max.Y)+grab_h_pixels, window_rect.Max.Y),
			},
		}
	}

	window.DrawList.AddRectFilledEx(grab_rect.Min, grab_rect.Max, grab_col, style.ScrollbarRounding, DrawCornerFlagsAll)
}

func (c *Context) GetBorderRect(window *Window, border_n int, perp_padding, thickness float64) f64.Rectangle {
	rect := window.Rect()
	if thickness == 0 {
		rect.Max = rect.Max.Sub(f64.Vec2{1, 1})
	}
	switch border_n {
	case 0:
		return f64.Rect(rect.Min.X+perp_padding, rect.Min.Y, rect.Max.X-perp_padding, rect.Min.Y+thickness)
	case 1:
		return f64.Rect(rect.Max.X-thickness, rect.Min.Y+perp_padding, rect.Max.X, rect.Max.Y-perp_padding)
	case 2:
		return f64.Rect(rect.Min.X+perp_padding, rect.Max.Y-thickness, rect.Max.X-perp_padding, rect.Max.Y)
	case 3:
		return f64.Rect(rect.Min.X, rect.Min.Y+perp_padding, rect.Min.X+thickness, rect.Max.Y-perp_padding)
	}

	return f64.Rectangle{}
}

func (c *Context) SetupDrawData(draw_lists []*DrawList, out_draw_data *DrawData) {
	out_draw_data.Valid = true
	out_draw_data.CmdLists = nil
	if len(draw_lists) > 0 {
		out_draw_data.CmdLists = draw_lists
	}
	out_draw_data.CmdListsCount = len(draw_lists)
	out_draw_data.TotalVtxCount, out_draw_data.TotalIdxCount = 0, 0
	for n := range draw_lists {
		out_draw_data.TotalVtxCount += len(draw_lists[n].VtxBuffer)
		out_draw_data.TotalIdxCount += len(draw_lists[n].IdxBuffer)
	}
}