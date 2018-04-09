package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type ComboFlags int

const (
	ComboFlagsPopupAlignLeft ComboFlags = 1 << 0 // Align the popup toward the left by default
	ComboFlagsHeightSmall    ComboFlags = 1 << 1 // Max ~4 items visible. Tip: If you want your combo popup to be a specific size you can use SetNextWindowSizeConstraints() prior to calling BeginCombo()
	ComboFlagsHeightRegular  ComboFlags = 1 << 2 // Max ~8 items visible (default)
	ComboFlagsHeightLarge    ComboFlags = 1 << 3 // Max ~20 items visible
	ComboFlagsHeightLargest  ComboFlags = 1 << 4 // As many fitting items as possible
	ComboFlagsNoArrowButton  ComboFlags = 1 << 5 // Display on the preview box without the square arrow button
	ComboFlagsNoPreview      ComboFlags = 1 << 6 // Display only a square arrow button
	ComboFlagsHeightMask_    ComboFlags = ComboFlagsHeightSmall | ComboFlagsHeightRegular | ComboFlagsHeightLarge | ComboFlagsHeightLargest
)

func (c *Context) BeginCombo(label, preview_value string) bool {
	return c.BeginComboEx(label, preview_value, 0)
}

func (c *Context) BeginComboEx(label, preview_value string, flags ComboFlags) bool {
	// Always consume the SetNextWindowSizeConstraint() call in our early return paths
	backup_next_window_size_constraint := c.NextWindowData.SizeConstraintCond
	c.NextWindowData.SizeConstraintCond = 0

	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	// Can't use both flags together
	assert((flags & (ComboFlagsNoArrowButton | ComboFlagsNoPreview)) != (ComboFlagsNoArrowButton | ComboFlagsNoPreview))

	style := c.Style
	id := window.GetID(label)

	arrow_size := 0.0
	if flags&ComboFlagsNoArrowButton == 0 {
		arrow_size = c.GetFrameHeight()
	}
	label_size := c.CalcTextSizeEx(label, true, -1)
	w := arrow_size
	if flags&ComboFlagsNoPreview == 0 {
		w = c.CalcItemWidth()
	}
	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{w, label_size.Y + style.FramePadding.Y*2.0}),
	}
	frame_bb_x := 0.0
	if label_size.X > 0.0 {
		frame_bb_x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(f64.Vec2{frame_bb_x, 0}),
	}
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	if !c.ItemAddEx(total_bb, id, &frame_bb) {
		return false
	}

	hovered, _, pressed := c.ButtonBehavior(frame_bb, id, 0)
	popup_open := c.IsPopupOpen(id)

	value_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Sub(f64.Vec2{arrow_size, 0.0}),
	}
	var frame_col color.RGBA
	if hovered {
		frame_col = c.GetColorFromStyle(ColFrameBgHovered)
	} else {
		frame_col = c.GetColorFromStyle(ColFrameBg)
	}
	c.RenderNavHighlight(frame_bb, id)
	if flags&ComboFlagsNoPreview == 0 {
		window.DrawList.AddRectFilledEx(
			frame_bb.Min,
			f64.Vec2{frame_bb.Max.X - arrow_size, frame_bb.Max.Y},
			frame_col, style.FrameRounding, DrawCornerFlagsLeft,
		)
	}
	if flags&ComboFlagsNoArrowButton == 0 {
		var col color.RGBA
		if popup_open || hovered {
			col = c.GetColorFromStyle(ColButtonHovered)
		} else {
			col = c.GetColorFromStyle(ColButton)
		}
		flags := DrawCornerFlagsRight
		if w <= arrow_size {
			flags = DrawCornerFlagsAll
		}

		window.DrawList.AddRectFilledEx(
			f64.Vec2{frame_bb.Max.X - arrow_size, frame_bb.Min.Y},
			frame_bb.Max,
			col, style.FrameRounding, flags,
		)
		c.RenderArrow(
			f64.Vec2{frame_bb.Max.X - arrow_size + style.FramePadding.Y, frame_bb.Min.Y + style.FramePadding.Y},
			DirDown,
		)
	}

	c.RenderFrameBorder(frame_bb.Min, frame_bb.Max, style.FrameRounding)
	if preview_value != "" && flags&ComboFlagsNoPreview == 0 {
		c.RenderTextClippedEx(frame_bb.Min.Add(style.FramePadding), value_bb.Max, preview_value, nil, f64.Vec2{0.0, 0.0}, nil)
	}
	if label_size.X > 0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, frame_bb.Min.Y + style.FramePadding.Y}, label)
	}

	if (pressed || c.NavActivateId == id) && !popup_open {
		if window.DC.NavLayerCurrent == 0 {
			window.NavLastIds[0] = id
		}
		c.OpenPopupEx(id)
		popup_open = true
	}

	if !popup_open {
		return false
	}

	if backup_next_window_size_constraint != 0 {
		c.NextWindowData.SizeConstraintCond = backup_next_window_size_constraint
		c.NextWindowData.SizeConstraintRect.Min.X = math.Max(c.NextWindowData.SizeConstraintRect.Min.X, w)
	} else {
		if flags&ComboFlagsHeightMask_ == 0 {
			flags |= ComboFlagsHeightRegular
		}
		// Only one
		assert(mathutil.IsPow2(int(flags & ComboFlagsHeightMask_)))
		popup_max_height_in_items := -1
		if flags&ComboFlagsHeightRegular != 0 {
			popup_max_height_in_items = 8
		} else if flags&ComboFlagsHeightSmall != 0 {
			popup_max_height_in_items = 4
		} else if flags&ComboFlagsHeightLarge != 0 {
			popup_max_height_in_items = 20
		}
		c.SetNextWindowSizeConstraints(
			f64.Vec2{w, 0.0},
			f64.Vec2{math.MaxFloat32, c.CalcMaxPopupHeightFromItemCount(popup_max_height_in_items)},
			nil,
		)
	}

	// Recycle windows based on depth
	name := fmt.Sprintf("##Combo_%02d", len(c.CurrentPopupStack))

	// Peak into expected window size so we can position it
	popup_window := c.FindWindowByName(name)
	if popup_window != nil && popup_window.WasActive {
		size_contents := c.CalcSizeContents(popup_window)
		size_expected := c.CalcSizeAfterConstraint(popup_window, c.CalcSizeAutoFit(popup_window, size_contents))
		if flags&ComboFlagsPopupAlignLeft != 0 {
			popup_window.AutoPosLastDirection = DirLeft
		}
		pos := c.FindBestWindowPosForPopupEx(frame_bb.BL(), size_expected, &popup_window.AutoPosLastDirection, frame_bb, PopupPositionPolicyComboBox)
		c.SetNextWindowPos(pos, 0, f64.Vec2{0, 0})
	}

	window_flags := WindowFlagsAlwaysAutoResize | WindowFlagsPopup | WindowFlagsNoTitleBar | WindowFlagsNoResize | WindowFlagsNoSavedSettings
	if !c.BeginEx(name, nil, window_flags) {
		// This should never happen as we tested for IsPopupOpen() above
		c.EndPopup()
		assert(false)
		return false
	}

	// Horizontally align ourselves with the framed text
	if style.FramePadding.X != style.WindowPadding.X {
		c.IndentEx(style.FramePadding.X - style.WindowPadding.X)
	}

	return true
}

func (c *Context) EndCombo() {
	style := &c.Style
	if style.FramePadding.X != style.WindowPadding.X {
		c.UnindentEx(style.FramePadding.X - style.WindowPadding.X)
	}
	c.EndPopup()
}
