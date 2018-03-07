package imgui

import (
	"image/color"

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
	_ = backup_hovered_window

	hovered = c.ItemHoverable(bb, id)

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if flags&ButtonFlagsPressedOnDragDropHold != 0 && c.DragDropActive && c.DragDropSourceFlags&DragDropFlagsSourceNoHoldToOpenOthers == 0 {
		if c.IsItemHovered(HoveredFlagsAllowWhenBlockedByActiveItem) {
			hovered = true
			c.SetHoveredID(id)
		}
	}

	return
}