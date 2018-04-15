package imgui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) BeginDragDropSource() bool {
	return false
}

func (c *Context) EndDragDropSource() {
}

func (c *Context) SetDragDropPayload(a interface{}) {
}

func (c *Context) DragInt(label string, v *int) bool {
	return c.DragIntEx(label, v, 1, 0, 0, "%.0f")
}

// NB: v_speed is float to allow adjusting the drag speed with more precision
func (c *Context) DragIntEx(label string, v *int, v_speed float64, v_min, v_max int, display_format string) bool {
	if display_format == "" {
		display_format = "%.0f"
	}
	v_f := float64(*v)
	value_changed := c.DragFloatEx(label, &v_f, v_speed, float64(v_min), float64(v_max), display_format, 1)
	*v = int(v_f)
	return value_changed
}

func (c *Context) DragIntN(label string, v []int, v_speed float64, v_min, v_max int, display_format string) bool {
	components := len(v)
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	value_changed := false
	c.BeginGroup()
	c.PushStringID(label)
	c.PushMultiItemsWidths(components)
	for i := 0; i < components; i++ {
		c.PushID(ID(i))
		if c.DragIntEx("##v", &v[i], v_speed, v_min, v_max, display_format) {
			value_changed = true
		}
		c.SameLineEx(0, c.Style.ItemInnerSpacing.X)
		c.PopID()
		c.PopItemWidth()
	}
	c.PopID()

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()

	return value_changed
}

func (c *Context) DragFloat(label string, v *float64) bool {
	return c.DragFloatEx(label, v, 1, 0, 0, "%.3f", 1)
}

func (c *Context) DragFloatEx(label string, v *float64, v_speed, v_min, v_max float64, display_format string, power float64) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := c.Style
	id := window.GetID(label)
	w := c.CalcItemWidth()

	label_size := c.CalcTextSizeEx(label, true, -1)
	frame_bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(f64.Vec2{w, label_size.Y + style.FramePadding.Y*2.0}),
	}
	inner_bb := f64.Rectangle{
		frame_bb.Min.Add(style.FramePadding),
		frame_bb.Max.Sub(style.FramePadding),
	}
	total_bb_x := 0.0
	if label_size.X > 0 {
		total_bb_x = style.ItemInnerSpacing.X + label_size.X
	}
	total_bb := f64.Rectangle{
		frame_bb.Min,
		frame_bb.Max.Add(f64.Vec2{total_bb_x, 0}),
	}

	// NB- we don't call ItemSize() yet because we may turn into a text edit box below
	if !c.ItemAddEx(total_bb, id, &frame_bb) {
		c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
		return false
	}
	hovered := c.ItemHoverable(frame_bb, id)

	if display_format == "" {
		display_format = "%.3f"
	}
	decimal_precision := ParseFormatPrecision(display_format, 3)

	// Tabbing or CTRL-clicking on Drag turns it into an input box
	start_text_input := false
	tab_focus_requested := c.FocusableItemRegister(window, id)
	if tab_focus_requested || (hovered && (c.IO.MouseClicked[0] || c.IO.MouseDoubleClicked[0])) || c.NavActivateId == id || (c.NavInputId == id && c.ScalarAsInputTextId != id) {
		c.SetActiveID(id, window)
		c.SetFocusID(id, window)
		c.FocusWindow(window)
		c.ActiveIdAllowNavDirFlags = (1 << uint(DirUp)) | (1 << uint(DirDown))
		if tab_focus_requested || c.IO.KeyCtrl || c.IO.MouseDoubleClicked[0] || c.NavInputId == id {
			start_text_input = true
			c.ScalarAsInputTextId = 0
		}
	}
	if start_text_input || (c.ActiveId == id && c.ScalarAsInputTextId == id) {
		return c.InputScalarAsWidgetReplacement(frame_bb, label, v, id, decimal_precision)
	}

	// Actual drag behavior
	c.ItemSizeBBEx(total_bb, style.FramePadding.Y)
	value_changed := c.DragBehavior(frame_bb, id, v, v_speed, v_min, v_max, decimal_precision, power)

	// Display value using user-provided display format so user can add prefix/suffix/decorations to the value.
	value := fmt.Sprintf(display_format, *v)
	c.RenderTextClippedEx(frame_bb.Min, frame_bb.Max, value, nil, f64.Vec2{0.5, 0.5}, nil)

	if label_size.X > 0.0 {
		c.RenderText(f64.Vec2{frame_bb.Max.X + style.ItemInnerSpacing.X, inner_bb.Min.Y}, label)
	}

	return value_changed
}

func (c *Context) DragFloatN(label string, v []float64, v_speed, v_min, v_max float64, display_format string, power float64) bool {
	components := len(v)
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	value_changed := false
	c.BeginGroup()
	c.PushStringID(label)
	c.PushMultiItemsWidths(components)
	for i := 0; i < components; i++ {
		c.PushID(ID(i))
		if c.DragFloatEx("##v", &v[i], v_speed, v_min, v_max, display_format, power) {
			value_changed = true
		}
		c.SameLineEx(0, c.Style.ItemInnerSpacing.X)
		c.PopID()
		c.PopItemWidth()
	}
	c.PopID()

	n := c.FindRenderedTextEnd(label)
	c.TextUnformatted(label[:n])
	c.EndGroup()

	return value_changed
}

func (c *Context) DragBehavior(frame_bb f64.Rectangle, id ID, v *float64, v_speed, v_min, v_max float64, decimal_precision int, power float64) bool {
	style := &c.Style

	// Draw frame
	var frame_col color.RGBA
	switch {
	case c.ActiveId == id:
		frame_col = c.GetColorFromStyle(ColFrameBgActive)
	case c.HoveredId == id:
		frame_col = c.GetColorFromStyle(ColFrameBgHovered)
	default:
		frame_col = c.GetColorFromStyle(ColFrameBg)
	}
	c.RenderNavHighlight(frame_bb, id)
	c.RenderFrameEx(frame_bb.Min, frame_bb.Max, frame_col, true, style.FrameRounding)

	value_changed := false

	// Process interacting with the drag
	if c.ActiveId == id {
		if c.ActiveIdSource == InputSourceMouse && !c.IO.MouseDown[0] {
			c.ClearActiveID()
		} else if c.ActiveIdSource == InputSourceNav && c.NavActivatePressedId == id && !c.ActiveIdIsJustActivated {
			c.ClearActiveID()
		}
	}
	if c.ActiveId == id {
		if c.ActiveIdIsJustActivated {
			// Lock current value on click
			c.DragCurrentValue = *v
			c.DragLastMouseDelta = f64.Vec2{0, 0}
		}

		if v_speed == 0.0 && (v_max-v_min) != 0.0 && (v_max-v_min) < math.MaxFloat32 {
			v_speed = (v_max - v_min) * c.DragSpeedDefaultRatio
		}

		v_cur := c.DragCurrentValue
		mouse_drag_delta := c.GetMouseDragDelta(0, 1.0)
		adjust_delta := 0.0
		if c.ActiveIdSource == InputSourceMouse && c.IsMousePosValid(nil) {
			adjust_delta = mouse_drag_delta.X - c.DragLastMouseDelta.X
			if c.IO.KeyShift && c.DragSpeedScaleFast >= 0.0 {
				adjust_delta *= c.DragSpeedScaleFast
			}
			if c.IO.KeyAlt && c.DragSpeedScaleSlow >= 0.0 {
				adjust_delta *= c.DragSpeedScaleSlow
			}
			c.DragLastMouseDelta.X = mouse_drag_delta.X
		}
		if c.ActiveIdSource == InputSourceNav {
			adjust_delta = c.GetNavInputAmount2dEx(NavDirSourceFlagsKeyboard|NavDirSourceFlagsPadDPad, InputReadModeRepeatFast, 1.0/10.0, 10.0).X
			// This is to avoid applying the saturation when already past the limits
			if v_min < v_max && ((v_cur >= v_max && adjust_delta > 0.0) || (v_cur <= v_min && adjust_delta < 0.0)) {
				adjust_delta = 0.0
			}
			v_speed = math.Max(v_speed, GetMinimumStepAtDecimalPrecision(decimal_precision))
		}
		adjust_delta *= v_speed

		if math.Abs(adjust_delta) > 0.0 {
			if math.Abs(power-1.0) > 0.001 {
				// Logarithmic curve on both side of 0.0
				v0_abs := math.Abs(v_cur)
				v0_sign := f64.SignNZ(v_cur)
				v1 := math.Pow(v0_abs, 1.0/power) + (adjust_delta * v0_sign)
				v1_abs := math.Abs(v1) // Crossed sign line
				v1_sign := f64.SignNZ(v1)
				v_cur = math.Pow(v1_abs, power) * v0_sign * v1_sign // Reapply sign
			} else {
				v_cur += adjust_delta
			}

			// Clamp
			if v_min < v_max {
				v_cur = f64.Clamp(v_cur, v_min, v_max)
			}
			c.DragCurrentValue = v_cur
		}

		// Round to user desired precision, then apply
		v_cur = f64.RoundPrec(v_cur, decimal_precision)
		if *v != v_cur {
			*v = v_cur
			value_changed = true
		}
	}

	return value_changed
}