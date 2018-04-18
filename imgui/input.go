package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type Key int

const (
	KeyTab Key = iota
	KeyLeftArrow
	KeyRightArrow
	KeyUpArrow
	KeyDownArrow
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyBackspace
	KeySpace
	KeyEnter
	KeyEscape
	KeyA // for text edit CTRL+A: select all
	KeyC // for text edit CTRL+C: copy
	KeyV // for text edit CTRL+V: paste
	KeyX // for text edit CTRL+X: cut
	KeyY // for text edit CTRL+Y: redo
	KeyZ // for text edit CTRL+Z: undo
	KeyCOUNT
)

type InputSource int

const (
	InputSourceNone InputSource = iota
	InputSourceMouse
	InputSourceNav
	InputSourceNavKeyboard // Only used occasionally for storage, not tested/handled by most code
	InputSourceNavGamepad  // "
	InputSourceCOUNT
)

type InputReadMode int

const (
	InputReadModeDown InputReadMode = iota
	InputReadModePressed
	InputReadModeReleased
	InputReadModeRepeat
	InputReadModeRepeatSlow
	InputReadModeRepeatFast
)

type MouseCursor int

const (
	MouseCursorNone MouseCursor = -1 + iota
	MouseCursorArrow
	MouseCursorTextInput  // When hovering over InputText, etc.
	MouseCursorResizeAll  // Unused
	MouseCursorResizeNS   // When hovering over an horizontal border
	MouseCursorResizeEW   // When hovering over a vertical border or a column
	MouseCursorResizeNESW // When hovering over the bottom-left corner of a window
	MouseCursorResizeNWSE // When hovering over the bottom-right corner of a window
	MouseCursorCOUNT
)

func (c *Context) CalcTypematicPressedRepeatAmount(t, t_prev, repeat_delay, repeat_rate float64) int {
	if t == 0 {
		return 1
	}
	if t <= repeat_delay || repeat_rate <= 0 {
		return 0
	}

	count := int((t-repeat_delay)/repeat_rate) - int((t_prev-repeat_delay)/repeat_rate)

	if count > 0 {
		return count
	}
	return 0
}

// Note that imgui doesn't know the semantic of each entry of io.KeysDown[]. Use your own indices/enums according to how your back-end/engine stored them into io.KeysDown[]!
func (c *Context) IsKeyDown(user_key_index int) bool {
	if user_key_index < 0 {
		return false
	}
	return c.IO.KeysDown[user_key_index]
}

func (c *Context) IsKeyPressed(user_key_index int, repeat bool) bool {
	if user_key_index < 0 {
		return false
	}
	t := c.IO.KeysDownDuration[user_key_index]
	if t == 0 {
		return true
	}
	if repeat && t > c.IO.KeyRepeatDelay {
		return c.GetKeyPressedAmount(user_key_index, c.IO.KeyRepeatDelay, c.IO.KeyRepeatRate) > 0
	}
	return false
}

func (c *Context) IsKeyReleased(user_key_index int) bool {
	if user_key_index < 0 {
		return false
	}
	return c.IO.KeysDownDurationPrev[user_key_index] >= 0 && !c.IO.KeysDown[user_key_index]
}

func (c *Context) IsMouseDown(button int) bool {
	return c.IO.MouseDown[button]
}

func (c *Context) GetKeyPressedAmount(key_index int, repeat_delay, repeat_rate float64) int {
	if key_index < 0 {
		return 0
	}
	t := c.IO.KeysDownDuration[key_index]
	return c.CalcTypematicPressedRepeatAmount(t, t-c.IO.DeltaTime, repeat_delay, repeat_rate)
}

func (c *Context) PushAllowKeyboardFocus(allow_keyboard_focus bool) {
	c.PushItemFlag(ItemFlagsAllowKeyboardFocus, allow_keyboard_focus)
}

func (c *Context) PopAllowKeyboardFocus() {
	c.PopItemFlag()
}

func (c *Context) IsKeyPressedMap(key Key) bool {
	return c.IsKeyPressedMapEx(key, true)
}

func (c *Context) IsKeyPressedMapEx(key Key, repeat bool) bool {
	key_index := c.IO.KeyMap[key]
	if key_index >= 0 {
		return c.IsKeyPressed(key_index, repeat)
	}
	return false
}

func (c *Context) InputInt(label string, v *int) bool {
	return c.InputIntEx(label, v, 1, 100, 0)
}

func (c *Context) InputIntEx(label string, v *int, step, step_fast int, extra_flags InputTextFlags) bool {
	// Hexadecimal input provided as a convenience but the flag name is awkward. Typically you'd use InputText() to parse your own data, if you want to handle prefixes.
	return false
}

func (c *Context) InputIntN(label string, v []int, extra_flags InputTextFlags) bool {
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
		if c.InputIntEx("##v", &v[i], 0, 0, extra_flags) {
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

func (c *Context) InputFloat(label string, v *float64, step float64) bool {
	return c.InputFloatEx(label, v, step, 0, -1, 0)
}

func (c *Context) InputFloatEx(label string, v *float64, step, step_fast float64, decimal_precision int, extra_flags InputTextFlags) bool {
	return false
}

// NB: scalar_format here must be a simple "%xx" format string with no prefix/suffix (unlike the Drag/Slider functions "display_format" argument)
func (c *Context) InputScalarEx(label string, data, step_ptr, step_fast_ptr interface{}, scalar_format string, extra_flags InputTextFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	label_size := c.CalcTextSizeEx(label, true, -1)

	c.BeginGroup()
	c.PushStringID(label)
	button_sz := f64.Vec2{c.GetFrameHeight(), c.GetFrameHeight()}
	if step_ptr != nil {
		c.PushItemWidth(math.Max(1.0, c.CalcItemWidth()-(button_sz.X+style.ItemInnerSpacing.X)*2))
	}

	buf := []byte(DataTypeFormatStringCustom(data, scalar_format))
	value_changed := false
	if (extra_flags & (InputTextFlagsCharsHexadecimal | InputTextFlagsCharsScientific)) == 0 {
		extra_flags |= InputTextFlagsCharsDecimal
	}
	extra_flags |= InputTextFlagsAutoSelectAll

	// PushId(label) + "" gives us the expected ID from outside point of view
	if c.InputTextEx("", buf, f64.Vec2{0, 0}, extra_flags, nil) {
		value_changed = DataTypeApplyOpFromText(buf, string(c.InputTextState.InitialText), data, scalar_format)
	}

	// Step buttons
	if step_ptr != nil {
		c.PopItemWidth()
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		if c.ButtonEx("-", button_sz, ButtonFlagsRepeat|ButtonFlagsDontClosePopups) {
			value_changed = true
		}
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		if c.ButtonEx("+", button_sz, ButtonFlagsRepeat|ButtonFlagsDontClosePopups) {
			value_changed = true
		}
	}
	c.PopID()

	if label_size.X > 0 {
		c.SameLineEx(0, style.ItemInnerSpacing.X)
		c.RenderText(f64.Vec2{window.DC.CursorPos.X, window.DC.CursorPos.Y + style.FramePadding.Y}, label)
		c.ItemSizeEx(label_size, style.FramePadding.Y)

	}
	c.EndGroup()

	return value_changed
}