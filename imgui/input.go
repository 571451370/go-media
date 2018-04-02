package imgui

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

// Note that imgui doesn't know the semantic of each entry of io.KeyDown[]. Use your own indices/enums according to how your back-end/engine stored them into KeyDown[]!
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
