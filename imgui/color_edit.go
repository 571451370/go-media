package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/math/mathutil"
)

type ColorEditFlags int

const (
	ColorEditFlagsNoAlpha        ColorEditFlags = 1 << 1 //              // ColorEdit ColorPicker ColorButton: ignore Alpha component (read 3 components from the input pointer).
	ColorEditFlagsNoPicker       ColorEditFlags = 1 << 2 //              // ColorEdit: disable picker when clicking on colored square.
	ColorEditFlagsNoOptions      ColorEditFlags = 1 << 3 //              // ColorEdit: disable toggling options menu when right-clicking on inputs/small preview.
	ColorEditFlagsNoSmallPreview ColorEditFlags = 1 << 4 //              // ColorEdit ColorPicker: disable colored square preview next to the inputs. (e.g. to show only the inputs)
	ColorEditFlagsNoInputs       ColorEditFlags = 1 << 5 //              // ColorEdit ColorPicker: disable inputs sliders/text widgets (e.g. to show only the small preview colored square).
	ColorEditFlagsNoTooltip      ColorEditFlags = 1 << 6 //              // ColorEdit ColorPicker ColorButton: disable tooltip when hovering the preview.
	ColorEditFlagsNoLabel        ColorEditFlags = 1 << 7 //              // ColorEdit ColorPicker: disable display of inline text label (the label is still forwarded to the tooltip and picker).
	ColorEditFlagsNoSidePreview  ColorEditFlags = 1 << 8 //              // ColorPicker: disable bigger color preview on right side of the picker use small colored square preview instead.

	// User Options (right-click on widget to change some of them). You can set application defaults using SetColorEditOptions(). The idea is that you probably don't want to override them in most of your calls let the user choose and/or call SetColorEditOptions() during startup.
	ColorEditFlagsAlphaBar         ColorEditFlags = 1 << 9  //              // ColorEdit ColorPicker: show vertical alpha bar/gradient in picker.
	ColorEditFlagsAlphaPreview     ColorEditFlags = 1 << 10 //              // ColorEdit ColorPicker ColorButton: display preview as a transparent color over a checkerboard instead of opaque.
	ColorEditFlagsAlphaPreviewHalf ColorEditFlags = 1 << 11 //              // ColorEdit ColorPicker ColorButton: display half opaque / half checkerboard instead of opaque.
	ColorEditFlagsHDR              ColorEditFlags = 1 << 12 //              // (WIP) ColorEdit: Currently only disable 0.0f..1.0f limits in RGBA edition (note: you probably want to use ColorEditFlagsFloat flag as well).
	ColorEditFlagsRGB              ColorEditFlags = 1 << 13 // [Inputs]     // ColorEdit: choose one among RGB/HSV/HEX. ColorPicker: choose any combination using RGB/HSV/HEX.
	ColorEditFlagsHSV              ColorEditFlags = 1 << 14 // [Inputs]     // "
	ColorEditFlagsHEX              ColorEditFlags = 1 << 15 // [Inputs]     // "
	ColorEditFlagsUint8            ColorEditFlags = 1 << 16 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0..255.
	ColorEditFlagsFloat            ColorEditFlags = 1 << 17 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0.0f..1.0f floats instead of 0..255 integers. No round-trip of value via integers.
	ColorEditFlagsPickerHueBar     ColorEditFlags = 1 << 18 // [PickerMode] // ColorPicker: bar for Hue rectangle for Sat/Value.
	ColorEditFlagsPickerHueWheel   ColorEditFlags = 1 << 19 // [PickerMode] // ColorPicker: wheel for Hue triangle for Sat/Value.

	// Internals/Masks
	ColorEditFlags_InputsMask     ColorEditFlags = ColorEditFlagsRGB | ColorEditFlagsHSV | ColorEditFlagsHEX
	ColorEditFlags_DataTypeMask   ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsFloat
	ColorEditFlags_PickerMask     ColorEditFlags = ColorEditFlagsPickerHueWheel | ColorEditFlagsPickerHueBar
	ColorEditFlags_OptionsDefault ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsRGB | ColorEditFlagsPickerHueBar // Change application default using SetColorEditOptions()
)

func (c *Context) SetColorEditOptions(flags ColorEditFlags) {
	if flags&ColorEditFlags_InputsMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_InputsMask
	}
	if flags&ColorEditFlags_DataTypeMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_DataTypeMask
	}
	if flags&ColorEditFlags_PickerMask == 0 {
		flags |= ColorEditFlags_OptionsDefault & ColorEditFlags_PickerMask
	}
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_InputsMask)))   // Check only 1 option is selected
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_DataTypeMask))) // Check only 1 option is selected
	assert(mathutil.IsPow2(int(flags & ColorEditFlags_PickerMask)))   // Check only 1 option is selected
	c.ColorEditOptions = flags
}

// Note: only access 3 floats if ImGuiColorEditFlags_NoAlpha flag is set.
func (c *Context) ColorTooltip(text string, col color.RGBA, flags ColorEditFlags) {
	cr, cg, cb, ca := col.R, col.G, col.B, col.A
	fc := chroma.RGBA2VEC4(col)

	c.BeginTooltipEx(0, true)

	text_end := c.FindRenderedTextEnd(text)
	if text_end > 0 {
		c.TextUnformatted(text[:text_end])
		c.Separator()
	}
	sz := f64.Vec2{
		c.FontSize*3 + c.Style.FramePadding.Y*2,
		c.FontSize*3 + c.Style.FramePadding.Y*2,
	}
	c.ColorButtonEx("##preview", col, flags&(ColorEditFlagsNoAlpha|ColorEditFlagsAlphaPreview|ColorEditFlagsAlphaPreviewHalf)|ColorEditFlagsNoTooltip, sz)
	c.SameLine()
	if flags&ColorEditFlagsNoAlpha != 0 {
		c.Text("#%02X%02X%02X\nR: %d, G: %d, B: %d\n(%.3f, %.3f, %.3f)", cr, cg, cb, cr, cg, cb, fc.X, fc.Y, fc.Z)
	} else {
		c.Text("#%02X%02X%02X%02X\nR:%d, G:%d, B:%d, A:%d\n(%.3f, %.3f, %.3f, %.3f)", cr, cg, cb, ca, cr, cg, cb, ca, fc.X, fc.Y, fc.Z, fc.W)
	}
	c.EndTooltip()
}
