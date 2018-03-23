package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type ButtonFlags int

const (
	ButtonFlagsRepeat                ButtonFlags = 1 << 0  // hold to repeat
	ButtonFlagsPressedOnClickRelease ButtonFlags = 1 << 1  // return true on click + release on same item [DEFAULT if no PressedOn* flag is set]
	ButtonFlagsPressedOnClick        ButtonFlags = 1 << 2  // return true on click (default requires click+release)
	ButtonFlagsPressedOnRelease      ButtonFlags = 1 << 3  // return true on release (default requires click+release)
	ButtonFlagsPressedOnDoubleClick  ButtonFlags = 1 << 4  // return true on double-click (default requires click+release)
	ButtonFlagsFlattenChildren       ButtonFlags = 1 << 5  // allow interactions even if a child window is overlapping
	ButtonFlagsAllowItemOverlap      ButtonFlags = 1 << 6  // require previous frame HoveredId to either match id or be null before being usable use along with SetItemAllowOverlap()
	ButtonFlagsDontClosePopups       ButtonFlags = 1 << 7  // disable automatically closing parent popup on press // [UNUSED]
	ButtonFlagsDisabled              ButtonFlags = 1 << 8  // disable interactions
	ButtonFlagsAlignTextBaseLine     ButtonFlags = 1 << 9  // vertically align button to match text baseline - ButtonEx() only // FIXME: Should be removed and handled by SmallButton() not possible currently because of DC.CursorPosPrevLine
	ButtonFlagsNoKeyModifiers        ButtonFlags = 1 << 10 // disable interaction if a key modifier is held
	ButtonFlagsNoHoldingActiveID     ButtonFlags = 1 << 11 // don't set ActiveId while holding the mouse (ButtonFlagsPressedOnClick only)
	ButtonFlagsPressedOnDragDropHold ButtonFlags = 1 << 12 // press when held into while we are drag and dropping another item (used by e.g. tree nodes collapsing headers)
	ButtonFlagsNoNavFocus            ButtonFlags = 1 << 13 // don't override navigation focus when activated
)

func (c *Context) ColorButton(desc_id string, col color.RGBA) bool {
	return c.ColorButtonEx(desc_id, col, 0, f64.Vec2{0, 0})
}

func (c *Context) ColorButtonEx(desc_id string, col color.RGBA, flags ColorEditFlags, size f64.Vec2) bool {
	return false
}

func (c *Context) Button(label string) bool {
	return c.ButtonEx(label, f64.Vec2{}, 0)
}

func (c *Context) SmallButton(label string) bool {
	backup_padding_y := c.Style.FramePadding.Y
	c.Style.FramePadding.Y = 0
	pressed := c.ButtonEx(label, f64.Vec2{0, 0}, ButtonFlagsAlignTextBaseLine)
	c.Style.FramePadding.Y = backup_padding_y
	return pressed
}

func (c *Context) ArrowButton(str_id string, dir Dir) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	id := window.GetID(str_id)
	sz := c.GetFrameHeight()
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{sz, sz}),
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, id) {
		return false
	}

	hovered, held, pressed := c.ButtonBehavior(bb, id, 0)

	var col color.RGBA
	switch {
	case hovered && held:
		col = c.GetColorFromStyle(ColButtonActive)
	case hovered:
		col = c.GetColorFromStyle(ColButtonHovered)
	default:
		col = c.GetColorFromStyle(ColButton)
	}

	// Render
	c.RenderNavHighlight(bb, id)
	c.RenderFrameEx(bb.Min, bb.Max, col, true, c.Style.FrameRounding)
	c.RenderArrow(bb.Min.Add(c.Style.FramePadding), dir)

	return pressed
}

func (c *Context) ButtonEx(label string, size_arg f64.Vec2, flags ButtonFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	pos := window.DC.CursorPos
	if flags&ButtonFlagsAlignTextBaseLine != 0 && style.FramePadding.Y < window.DC.CurrentLineTextBaseOffset {
		pos.Y += window.DC.CurrentLineTextBaseOffset - style.FramePadding.Y
	}
	size := c.CalcItemSize(size_arg, label_size.X+style.FramePadding.X*2, label_size.Y+style.FramePadding.Y*2)

	bb := f64.Rectangle{pos, pos.Add(size)}
	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	if !c.ItemAdd(bb, id) {
		return false
	}

	if window.DC.ItemFlags&ItemFlagsButtonRepeat != 0 {
		flags |= ButtonFlagsRepeat
	}

	hovered, held, pressed := c.ButtonBehavior(bb, id, flags)

	var col color.RGBA
	switch {
	case hovered && held:
		col = c.GetColorFromStyle(ColButtonActive)
	case hovered:
		col = c.GetColorFromStyle(ColButtonHovered)
	default:
		col = c.GetColorFromStyle(ColButton)
	}

	// Render
	c.RenderNavHighlight(bb, id)
	c.RenderFrameEx(bb.Min, bb.Max, col, true, style.FrameRounding)
	c.RenderTextClippedEx(
		bb.Min.Add(style.FramePadding),
		bb.Max.Sub(style.FramePadding),
		label,
		&label_size,
		style.ButtonTextAlign,
		&bb,
	)

	return pressed
}

func (c *Context) ButtonBehavior(bb f64.Rectangle, id ID, flags ButtonFlags) (hovered, held, pressed bool) {
	window := c.GetCurrentWindow()

	if flags&ButtonFlagsDisabled != 0 {
		if c.ActiveId == id {
			c.ClearActiveID()
		}
		return
	}

	// Default behavior requires click+release on same spot
	if flags&(ButtonFlagsPressedOnClickRelease|ButtonFlagsPressedOnClick|ButtonFlagsPressedOnRelease|ButtonFlagsPressedOnDoubleClick) == 0 {
		flags |= ButtonFlagsPressedOnClickRelease
	}

	backup_hovered_window := c.HoveredWindow
	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = window
	}

	hovered = c.ItemHoverable(bb, id)

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if flags&ButtonFlagsPressedOnDragDropHold != 0 && c.DragDropActive && c.DragDropSourceFlags&DragDropFlagsSourceNoHoldToOpenOthers == 0 {
		if c.IsItemHovered(HoveredFlagsAllowWhenBlockedByActiveItem) {
			hovered = true
			c.SetHoveredID(id)
			// FIXME: Our formula for CalcTypematicPressedRepeatAmount() is fishy
			if c.CalcTypematicPressedRepeatAmount(c.HoveredIdTimer+0.0001, c.HoveredIdTimer+0.0001-c.IO.DeltaTime, 0.01, 0.70) != 0 {
				pressed = true
				c.FocusWindow(window)
			}
		}
	}

	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = backup_hovered_window
	}

	// AllowOverlap mode (rarely used) requires previous frame HoveredId to be null or to match. This allows using patterns where a later submitted widget overlaps a previous one.
	if hovered && flags&ButtonFlagsAllowItemOverlap != 0 && (c.HoveredIdPreviousFrame != id && c.HoveredIdPreviousFrame != 0) {
		hovered = false
	}

	// Mouse
	if hovered {
		if flags&ButtonFlagsNoKeyModifiers == 0 || (!c.IO.KeyCtrl && !c.IO.KeyShift && !c.IO.KeyAlt) {
			//                        | CLICKING        | HOLDING with ImGuiButtonFlags_Repeat
			// PressedOnClickRelease  |  <on release>*  |  <on repeat> <on repeat> .. (NOT on release)  <-- MOST COMMON! (*) only if both click/release were over bounds
			// PressedOnClick         |  <on click>     |  <on click> <on repeat> <on repeat> ..
			// PressedOnRelease       |  <on release>   |  <on repeat> <on repeat> .. (NOT on release)
			// PressedOnDoubleClick   |  <on dclick>    |  <on dclick> <on repeat> <on repeat> ..
			// FIXME-NAV: We don't honor those different behaviors.
			if flags&ButtonFlagsPressedOnClickRelease != 0 && c.IO.MouseClicked[0] {
				c.SetActiveID(id, window)
				if flags&ButtonFlagsNoNavFocus == 0 {
					c.SetFocusID(id, window)
				}
				c.FocusWindow(window)
			}

			if (flags&ButtonFlagsPressedOnClick != 0 && c.IO.MouseClicked[0]) || (flags&ButtonFlagsPressedOnDoubleClick != 0 && c.IO.MouseDoubleClicked[0]) {
				pressed = true
				if flags&ButtonFlagsNoHoldingActiveID != 0 {
					c.ClearActiveID()
				} else {
					c.SetActiveID(id, window) // Hold on ID
				}
				c.FocusWindow(window)
			}

			if flags&ButtonFlagsPressedOnRelease != 0 && c.IO.MouseReleased[0] {
				// Repeat mode trumps <on release>
				if !(flags&ButtonFlagsRepeat == 0 && c.IO.MouseDownDurationPrev[0] >= c.IO.KeyRepeatDelay) {
					pressed = true
				}
				c.ClearActiveID()
			}

			// 'Repeat' mode acts when held regardless of _PressedOn flags (see table above).
			// Relies on repeat logic of IsMouseClicked() but we may as well do it ourselves if we end up exposing finer RepeatDelay/RepeatRate settings.
			if flags&ButtonFlagsRepeat != 0 && c.ActiveId == id && c.IO.MouseDownDuration[0] > 0 && c.IsMouseClicked(0, true) {
				pressed = true
			}
		}

		if pressed {
			c.NavDisableHighlight = true
		}
	}

	// Gamepad/Keyboard navigation
	// We report navigated item as hovered but we don't set g.HoveredId to not interfere with mouse.
	if c.NavId == id && !c.NavDisableHighlight && c.NavDisableMouseHover && (c.ActiveId == 0 || c.ActiveId == id || c.ActiveId == window.MoveId) {
		hovered = true
	}

	if c.NavActivateDownId == id {
		var nav_activated_by_inputs bool
		nav_activated_by_code := c.NavActivateId == id

		if flags&ButtonFlagsRepeat != 0 {
			nav_activated_by_inputs = c.IsNavInputPressed(NavInputActivate, InputReadModeRepeat)
		} else {
			nav_activated_by_inputs = c.IsNavInputPressed(NavInputActivate, InputReadModePressed)
		}

		if nav_activated_by_code || nav_activated_by_inputs {
			pressed = true
		}

		if nav_activated_by_code || nav_activated_by_inputs || c.ActiveId == id {
			// Set active id so it can be queried by user via IsItemActive(), equivalent of holding the mouse button.
			c.NavActivateId = id
			c.SetActiveID(id, window)
			if flags&ButtonFlagsNoNavFocus == 0 {
				c.SetFocusID(id, window)
			}
			c.ActiveIdAllowNavDirFlags = 1<<uint(DirLeft) | 1<<uint(DirRight) | 1<<uint(DirUp) | 1<<uint(DirDown)
		}
	}

	if c.ActiveId == id {
		if c.ActiveIdSource == InputSourceMouse {
			if c.ActiveIdIsJustActivated {
				c.ActiveIdClickOffset = c.IO.MousePos.Sub(bb.Min)
			}

			if c.IO.MouseDown[0] {
				held = true
			} else {
				if hovered && flags&ButtonFlagsPressedOnClickRelease != 0 {
					// Repeat mode trumps <on release>
					if !(flags&ButtonFlagsRepeat != 0 && c.IO.MouseDownDurationPrev[0] >= c.IO.KeyRepeatDelay) {
						if !c.DragDropActive {
							pressed = true
						}
					}
				}

				c.ClearActiveID()
			}

			if flags&ButtonFlagsNoNavFocus == 0 {
				c.NavDisableHighlight = true
			}
		} else if c.ActiveIdSource == InputSourceNav {
			if c.NavActivateDownId != id {
				c.ClearActiveID()
			}
		}
	}

	return
}

// Tip: use ImGui::PushID()/PopID() to push indices or pointers in the ID stack.
// Then you can keep 'str_id' empty or the same for all your buttons (instead of creating a string based on a non-string id)
func (c *Context) InvisibleButton(str_id string, size_arg f64.Vec2) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	id := window.GetID(str_id)
	size := c.CalcItemSize(size_arg, 0, 0)
	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(size),
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, id) {
		return false
	}

	_, _, pressed := c.ButtonBehavior(bb, id, 0)
	return pressed
}

// Button to close a window
func (c *Context) CloseButton(id ID, pos f64.Vec2, radius float64) bool {
	window := c.CurrentWindow

	// We intentionally allow interaction when clipped so that a mechanical Alt,Right,Validate sequence close a window.
	// (this isn't the regular behavior of buttons, but it doesn't affect the user much because navigation tends to keep items visible).
	rad := f64.Vec2{radius, radius}
	bb := f64.Rectangle{pos.Sub(rad), pos.Add(rad)}
	is_clipped := !c.ItemAdd(bb, id)

	hovered, held, pressed := c.ButtonBehavior(bb, id, 0)
	if is_clipped {
		return pressed
	}

	// Render
	center := bb.Center()
	if hovered {
		var col color.RGBA
		switch {
		case held && hovered:
			col = c.GetColorFromStyle(ColButtonActive)
		default:
			col = c.GetColorFromStyle(ColButtonHovered)
		}
		window.DrawList.AddCircleFilledEx(center, math.Max(2, radius), col, 9)
	}

	cross_extent := (radius * 0.7071) - 1.0
	cross_col := c.GetColorFromStyle(ColText)
	center = center.Sub(f64.Vec2{0.5, 0.5})
	window.DrawList.AddLineEx(center.Add(f64.Vec2{+cross_extent, +cross_extent}), center.Add(f64.Vec2{-cross_extent, -cross_extent}), cross_col, 1.0)
	window.DrawList.AddLineEx(center.Add(f64.Vec2{+cross_extent, -cross_extent}), center.Add(f64.Vec2{-cross_extent, +cross_extent}), cross_col, 1.0)

	return pressed
}

func (c *Context) PushButtonRepeat(repeat bool) {
	c.PushItemFlag(ItemFlagsButtonRepeat, repeat)
}

func (c *Context) PopButtonRepeat() {
	c.PopItemFlag()
}