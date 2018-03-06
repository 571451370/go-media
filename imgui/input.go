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