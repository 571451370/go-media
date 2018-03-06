package imgui

import "github.com/qeedquan/go-media/math/f64"

type NavInput int

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
	NavInputKeyMenu_  // toggle menu                                  // = io.KeyAlt
	NavInputKeyLeft_  // move left                                    // = Arrow keys
	NavInputKeyRight_ // move right
	NavInputKeyUp_    // move up
	NavInputKeyDown_  // move down
	NavInputCOUNT
	NavInputInternalStart_ = NavInputKeyMenu_
)

type Dir int

const (
	DirNone  Dir = -1
	DirLeft  Dir = 0
	DirRight Dir = 1
	DirUp    Dir = 2
	DirDown  Dir = 3
	DirCOUNT Dir = 4
)

type Cond int

const (
	CondAlways       Cond = 1 << 0 // Set the variable
	CondOnce         Cond = 1 << 1 // Set the variable once per runtime session (only the first call with succeed)
	CondFirstUseEver Cond = 1 << 2 // Set the variable if the window has no saved data (if doesn't exist in the .ini file)
	CondAppearing    Cond = 1 << 3 // Set the variable if the window is appearing after being hidden/inactive (or the first time)
)

type NavHighlightFlags int

const (
	NavHighlightFlagsTypeDefault NavHighlightFlags = 1 << 0
	NavHighlightFlagsTypeThin    NavHighlightFlags = 1 << 1
	NavHighlightFlagsAlwaysDraw  NavHighlightFlags = 1 << 2
	NavHighlightFlagsNoRounding  NavHighlightFlags = 1 << 3
)

type NavDirSourceFlags int

const (
	NavDirSourceFlagsKeyboard  NavDirSourceFlags = 1 << 0
	NavDirSourceFlagsPadDPad   NavDirSourceFlags = 1 << 1
	NavDirSourceFlagsPadLStick NavDirSourceFlags = 1 << 2
)

type NavForward int

const (
	NavForwardNone NavForward = iota
	NavForwardForwardQueued
	NavForwardForwardActive
)

type NavMoveResult struct {
	ID         ID      // Best candidate
	ParentID   ID      // Best candidate window->IDStack.back() - to compare context
	Window     *Window // Best candidate window
	DistBox    float64 // Best candidate box distance to current NavId
	DistCenter float64 // Best candidate center distance to current NavId
	DistAxial  float64
	RectRel    f64.Rectangle // Best candidate bounding box in window relative space
}