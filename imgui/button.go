package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/math/f64"
)

type ButtonFlags uint

const (
	ButtonFlagsRepeat ButtonFlags = 1 << iota
	ButtonFlagsPressedOnClickRelease
	ButtonFlagsPressedOnClick
	ButtonFlagsPressedOnRelease
	ButtonFlagsPressedOnDoubleClick
	ButtonFlagsFlattenChildren
	ButtonFlagsAllowItemOverlap
	ButtonFlagsDontClosePopups
	ButtonFlagsDisabled
	ButtonFlagsAlignTextBaseLine
	ButtonFlagsNoKeyModifiers
	ButtonFlagsNoHoldingActiveID
	ButtonFlagsPressedOnDragDropHold
	ButtonFlagsNoNavFocus
)

func (c *Context) Button(label string, size f64.Vec2) bool {
	return c.ButtonEx(label, size, 0)
}

func (c *Context) SmallButton(label string) bool {
	style := c.Style
	y := style.FramePadding.Y
	style.FramePadding.Y = 0
	pressed := c.ButtonEx(label, f64.Vec2{}, ButtonFlagsAlignTextBaseLine)
	style.FramePadding.Y = y
	return pressed
}

func (c *Context) ArrowButton(strId string, dir Dir) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	dc := &window.DC
	id := window.GetID(strId)
	sz := c.GetFrameHeight()
	bb := f64.Rectangle{dc.CursorPos, dc.CursorPos.Add(f64.Vec2{sz, sz})}
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

	// render
	c.RenderNavHighlight(bb, id)
	c.RenderFrame(bb.Min, bb.Max, col, true, style.FrameRounding)
	c.RenderArrow(bb.Min.Add(style.FramePadding), dir)

	return pressed
}

func (c *Context) ButtonEx(label string, size f64.Vec2, flags ButtonFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	dc := &window.DC
	id := window.GetID(label)
	labelSize := c.CalcTextSize(label, true, -1)

	pos := dc.CursorPos
	// Try to vertically align buttons that are smaller/have no padding so that text baseline matches
	// (bit hacky, since it shouldn't be a flag)
	if flags&ButtonFlagsAlignTextBaseLine != 0 && style.FramePadding.Y < dc.CurrentLineTextBaseOffset {
		pos.Y += dc.CurrentLineTextBaseOffset - style.FramePadding.Y
	}

	sz := c.CalcItemSize(size, labelSize.X+style.FramePadding.X*2, labelSize.Y+style.FramePadding.Y*2)
	bb := f64.Rectangle{pos, pos.Add(sz)}
	c.ItemSizeBBEx(bb, style.FramePadding.Y)
	if !c.ItemAdd(bb, id) {
		return false
	}

	if dc.ItemFlags&ItemFlagsButtonRepeat != 0 {
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

	// render
	c.RenderNavHighlight(bb, id)
	c.RenderFrame(bb.Min, bb.Max, col, true, style.FrameRounding)
	c.RenderTextClippedEx(bb.Min.Add(style.FramePadding), bb.Max.Sub(style.FramePadding), label, &labelSize, style.ButtonTextAlign, &bb)

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
	mask := ButtonFlagsPressedOnClickRelease | ButtonFlagsPressedOnClick |
		ButtonFlagsPressedOnRelease | ButtonFlagsPressedOnDoubleClick
	if flags&mask == 0 {
		flags |= ButtonFlagsPressedOnClickRelease
	}

	backupHoveredWindow := c.HoveredWindow
	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = window
	}
	hovered = c.ItemHoverable(bb, id)

	// Special mode for Drag and Drop where holding button pressed for a long time while dragging another item triggers the button
	if flags&ButtonFlagsPressedOnDragDropHold != 0 && c.DragDropActive && c.DragDropSourceFlags&DragDropFlagsSourceNoHoldToOpenOthers == 0 {
	}

	if flags&ButtonFlagsFlattenChildren != 0 && c.HoveredRootWindow == window {
		c.HoveredWindow = backupHoveredWindow
	}

	return
}
