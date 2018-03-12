package imgui

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type SliderFlags int

const (
	SliderFlagsVertical SliderFlags = 1 << 0
)

func (c *Context) SliderBehaviorCalcRatioFromValue(v, v_min, v_max, power, linear_zero_pos float64) float64 {
	if v_min == v_max {
		return 0
	}

	is_non_linear := (power < 1.0-0.00001) || (power > 1.0+0.00001)
	var v_clamped float64
	if v_min < v_max {
		v_clamped = f64.Clamp(v, v_min, v_max)
	} else {
		v_clamped = f64.Clamp(v, v_max, v_min)
	}

	if is_non_linear {
		if v_clamped < 0 {
			f := 1.0 - (v_clamped-v_min)/(math.Min(0.0, v_max)-v_min)
			return (1.0 - math.Pow(f, 1.0/power)) * linear_zero_pos
		} else {
			f := (v_clamped - math.Max(0.0, v_min)) / (v_max - math.Max(0.0, v_min))
			return linear_zero_pos + math.Pow(f, 1.0/power)*(1.0-linear_zero_pos)
		}
	}

	// Linear slider
	return (v_clamped - v_min) / (v_max - v_min)
}

func (c *Context) SliderBehavior(frame_bb f64.Rectangle, id ID, v *float64, v_min, v_max, power float64, decimal_precision int, flags SliderFlags) bool {
	window := c.GetCurrentWindow()
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

	is_non_linear := (power < 1.0-0.00001) || (power > 1.0+0.00001)
	is_horizontal := (flags & SliderFlagsVertical) == 0

	grab_padding := 2.0
	var slider_sz float64
	if is_horizontal {
		slider_sz = frame_bb.Dx() - grab_padding*2.0
	} else {
		slider_sz = frame_bb.Dy() - grab_padding*2.0
	}

	grab_sz := 0.0
	if decimal_precision != 0 {
		grab_sz = math.Min(style.GrabMinSize, slider_sz)
	} else {
		// Integer sliders, if possible have the grab size represent 1 unit
		v_ratio := math.Abs(v_max-v_min) + 1
		grab_sz = math.Min(math.Max(slider_sz/v_ratio, style.GrabMinSize), slider_sz)
	}
	slider_usable_sz := slider_sz - grab_sz
	var slider_usable_pos_min, slider_usable_pos_max float64
	if is_horizontal {
		slider_usable_pos_min = frame_bb.Min.X + grab_padding + grab_sz*0.5
		slider_usable_pos_max = frame_bb.Max.X - grab_padding - grab_sz*0.5
	} else {
		slider_usable_pos_min = frame_bb.Min.Y + grab_padding + grab_sz*0.5
		slider_usable_pos_max = frame_bb.Max.Y - grab_padding - grab_sz*0.5
	}

	// For logarithmic sliders that cross over sign boundary we want the exponential increase to be symmetric around 0.0f
	linear_zero_pos := 0.0 // 0.0->1.0f
	if v_min*v_max < 0.0 {
		// Different sign
		linear_dist_min_to_0 := math.Pow(math.Abs(0.0-v_min), 1.0/power)
		linear_dist_max_to_0 := math.Pow(math.Abs(v_max-0.0), 1.0/power)
		linear_zero_pos = linear_dist_min_to_0 / (linear_dist_min_to_0 + linear_dist_max_to_0)
	} else {
		// Same sign
		linear_zero_pos = 0
		if v_min < 0 {
			linear_zero_pos = 1
		}
	}

	// Process interacting with the slider
	value_changed := false
	if c.ActiveId == id {
		set_new_value := false
		clicked_t := 0.0
		if c.ActiveIdSource == InputSourceMouse {
			if !c.IO.MouseDown[0] {
				c.ClearActiveID()
			} else {
				mouse_abs_pos := c.IO.MousePos.Y
				if is_horizontal {
					mouse_abs_pos = c.IO.MousePos.X
				}
				clicked_t := 0.0
				if slider_usable_sz > 0 {
					clicked_t = f64.Clamp((mouse_abs_pos-slider_usable_pos_min)/slider_usable_sz, 0.0, 1.0)
				}
				if !is_horizontal {
					clicked_t = 1.0 - clicked_t
				}
				set_new_value = true
			}
		} else if c.ActiveIdSource == InputSourceNav {
			delta2 := c.GetNavInputAmount2d(NavDirSourceFlagsKeyboard|NavDirSourceFlagsPadDPad, InputReadModeRepeatFast, 0.0, 0.0)
			delta := -delta2.Y
			if is_horizontal {
				delta = delta2.X
			}
			if c.NavActivatePressedId == id && !c.ActiveIdIsJustActivated {
				c.ClearActiveID()
			} else if delta != 0.0 {
				clicked_t = c.SliderBehaviorCalcRatioFromValue(*v, v_min, v_max, power, linear_zero_pos)
				if decimal_precision == 0 && !is_non_linear {
					if math.Abs(v_max-v_min) <= 100.0 || c.IsNavInputDown(NavInputTweakSlow) {
						// Gamepad/keyboard tweak speeds in integer steps
						if delta < 0 {
							delta = -1 / (v_max - v_min)
						} else {
							delta = 1 / (v_max - v_min)
						}
					} else {
						delta /= 100.0
					}
				} else {
					delta /= 100.0 // Gamepad/keyboard tweak speeds in % of slider bounds
					if c.IsNavInputDown(NavInputTweakSlow) {
						delta /= 10.0
					}
				}

				if c.IsNavInputDown(NavInputTweakFast) {
					delta *= 10.0
				}
				set_new_value = true
				// This is to avoid applying the saturation when already past the limits
				if (clicked_t >= 1.0 && delta > 0.0) || (clicked_t <= 0.0 && delta < 0.0) {
					set_new_value = false
				} else {
					clicked_t = f64.Saturate(clicked_t + delta)
				}
			}
		}

		if set_new_value {
			var new_value float64

			if is_non_linear {
				// Account for logarithmic scale on both sides of the zero
				if clicked_t < linear_zero_pos {
					// Negative: rescale to the negative range before powering
					a := 1.0 - (clicked_t / linear_zero_pos)
					a = math.Pow(a, power)
					new_value = f64.Lerp(a, math.Min(v_max, 0.0), v_min)
				} else {
					// Positive: rescale to the positive range before powering
					var a float64
					if math.Abs(linear_zero_pos-1.0) > 1.e-6 {
						a = (clicked_t - linear_zero_pos) / (1.0 - linear_zero_pos)
					} else {
						a = clicked_t
					}
					a = math.Pow(a, power)
					new_value = f64.Lerp(a, math.Max(v_min, 0.0), v_max)
				}
			} else {
				// Linear slider
				new_value = f64.Lerp(clicked_t, v_min, v_max)
			}

			// Round past decimal precision
			new_value = f64.RoundPrec(new_value, decimal_precision)
			if *v != new_value {
				*v = new_value
				value_changed = true
			}
		}
	}

	grab_t := c.SliderBehaviorCalcRatioFromValue(*v, v_min, v_max, power, linear_zero_pos)
	if !is_horizontal {
		grab_t = 1.0 - grab_t
	}
	grab_pos := f64.Lerp(grab_t, slider_usable_pos_min, slider_usable_pos_max)
	var grab_bb f64.Rectangle
	if is_horizontal {
		grab_bb = f64.Rectangle{
			f64.Vec2{grab_pos - grab_sz*0.5, frame_bb.Min.Y + grab_padding},
			f64.Vec2{grab_pos + grab_sz*0.5, frame_bb.Max.Y - grab_padding},
		}
	} else {
		grab_bb = f64.Rectangle{
			f64.Vec2{frame_bb.Min.X + grab_padding, grab_pos - grab_sz*0.5},
			f64.Vec2{frame_bb.Max.X - grab_padding, grab_pos + grab_sz*0.5},
		}
	}

	var col color.RGBA
	if c.ActiveId == id {
		col = c.GetColorFromStyle(ColSliderGrabActive)
	} else {
		col = c.GetColorFromStyle(ColSliderGrabActive)
	}
	window.DrawList.AddRectFilledEx(grab_bb.Min, grab_bb.Max, col, style.GrabRounding, DrawCornerFlagsAll)

	return value_changed
}