package imgui

type Key uint

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
	KeyCount
)

type NavInput uint

const (
	// Gamepad Mapping
	NavInputActivate    NavInput = iota // activate / open / toggle / tweak value       // e.g. Circle (PS4) A (Xbox) A (Switch) Space (Keyboard)
	NavInputCancel                      // cancel / close / exit                        // e.g. Cross  (PS4) B (Xbox) B (Switch) Escape (Keyboard)
	NavInputInput                       // text input / on-screen keyboard              // e.g. Triang.(PS4) Y (Xbox) X (Switch) Return (Keyboard)
	NavInputMenu                        // tap: toggle menu / hold: focus move resize // e.g. Square (PS4) X (Xbox) Y (Switch) Alt (Keyboard)
	NavInputDpadLeft                    // move / tweak / resize window (w/ PadMenu)    // e.g. D-pad Left/Right/Up/Down (Gamepads) Arrow keys (Keyboard)
	NavInputDpadRight                   //
	NavInputDpadUp                      //
	NavInputDpadDown                    //
	NavInputLStickLeft                  // scroll / move window (w/ PadMenu)            // e.g. Left Analog Stick Left/Right/Up/Down
	NavInputLStickRight                 //
	NavInputLStickUp                    //
	NavInputLStickDown                  //
	NavInputFocusPrev                   // next window (w/ PadMenu)                     // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputFocusNext                   // prev window (w/ PadMenu)                     // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)
	NavInputTweakSlow                   // slower tweaks                                // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputTweakFast                   // faster tweaks                                // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)

	// [Internal] Don't use directly! This is used internally to differentiate keyboard from gamepad inputs for behaviors that require to differentiate them.
	// Keyboard behavior that have no corresponding gamepad mapping (e.g. CTRL+TAB) may be directly reading from io.KeyDown[] instead of io.NavInputs[].
	NavInputKeyMenu  // toggle menu                                  // = io.KeyAlt
	NavInputKeyLeft  // move left                                    // = Arrow keys
	NavInputKeyRight // move right
	NavInputKeyUp    // move up
	NavInputKeyDown  // move down
	NavInputCount
	NavInputInternalStart = NavInputKeyMenu
)

type InputSource uint

const (
	InputSourceNone InputSource = iota
	InputSourceMouse
	InputSourceNav
	InputSourceNavKeyboard
	InputSourceNavGamepad
	InputSourceCount
)

type InputReadMode uint

const (
	InputReadModeDown InputReadMode = iota
	InputReadModePressed
	InputReadModeReleased
	InputReadModeRepeat
	InputReadModeRepeatSlow
	InputReadModeRepeatFast
)
