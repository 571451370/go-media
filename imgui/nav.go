package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

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

func (c *Context) IsNavInputPressed(n NavInput, mode InputReadMode) bool {
	return c.GetNavInputAmount(n, mode) > 0
}

func (c *Context) GetNavInputAmount(n NavInput, mode InputReadMode) float64 {
	// Instant, read analog input (0.0f..1.0f, as provided by user)
	if mode == InputReadModeDown {
		return c.IO.NavInputs[n]
	}

	t := c.IO.NavInputsDownDuration[n]
	// Return 1.0f when just released, no repeat, ignore analog input.
	if t < 0 && mode == InputReadModeReleased {
		if c.IO.NavInputsDownDurationPrev[n] >= 0 {
			return 1
		}
		return 0
	}

	if t < 0 {
		return 0
	}

	// Return 1.0f when just pressed, no repeat, ignore analog input.
	if mode == InputReadModePressed {
		if t == 0 {
			return 1
		}
		return 0
	}

	if mode == InputReadModeRepeat {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.80,
		))
	}

	if mode == InputReadModeRepeatSlow {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*1,
			c.IO.KeyRepeatRate*2,
		))
	}

	if mode == InputReadModeRepeatFast {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.30,
		))
	}

	return 0
}

func (c *Context) GetNavInputAmount2d(dir_sources NavDirSourceFlags, mode InputReadMode, slow_factor, fast_factor float64) f64.Vec2 {
	delta := f64.Vec2{}
	if dir_sources&NavDirSourceFlagsKeyboard != 0 {
		right := c.GetNavInputAmount(NavInputKeyRight_, mode)
		left := c.GetNavInputAmount(NavInputKeyLeft_, mode)
		down := c.GetNavInputAmount(NavInputKeyRight_, mode)
		up := c.GetNavInputAmount(NavInputKeyRight_, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadDPad != 0 {
		right := c.GetNavInputAmount(NavInputDpadRight, mode)
		left := c.GetNavInputAmount(NavInputDpadLeft, mode)
		down := c.GetNavInputAmount(NavInputDpadDown, mode)
		up := c.GetNavInputAmount(NavInputDpadUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadLStick != 0 {
		right := c.GetNavInputAmount(NavInputLStickRight, mode)
		left := c.GetNavInputAmount(NavInputLStickLeft, mode)
		down := c.GetNavInputAmount(NavInputLStickDown, mode)
		up := c.GetNavInputAmount(NavInputLStickUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if slow_factor != 0.0 && c.IsNavInputDown(NavInputTweakSlow) {
		delta = delta.Scale(slow_factor)
	}
	if fast_factor != 0.0 && c.IsNavInputDown(NavInputTweakFast) {
		delta = delta.Scale(fast_factor)
	}
	return delta
}

// FIXME-OPT O(N)
func (c *Context) FindWindowIndex(window *Window) int {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		if c.Windows[i] == window {
			return i
		}
	}
	return -1
}

// FIXME-OPT O(N)
func (c *Context) FindWindowNavigable(i_start, i_stop, dir int) *Window {
	for i := i_start; i >= 0 && i < len(c.Windows) && i != i_stop; i += dir {
		if c.IsWindowNavFocusable(c.Windows[i]) {
			return c.Windows[i]
		}
	}
	return nil
}

func (c *Context) IsWindowNavFocusable(window *Window) bool {
	return window.Active && window == window.RootWindowForTabbing && window.Flags&WindowFlagsNoNavFocus == 0 || window == c.NavWindow
}

func (c *Context) NavUpdateWindowingHighlightWindow(focus_change_dir int) {
	if c.NavWindowingTarget.Flags&WindowFlagsModal != 0 {
		return
	}

	i_current := c.FindWindowIndex(c.NavWindowingTarget)
	window_target := c.FindWindowNavigable(i_current+focus_change_dir, -math.MinInt32, focus_change_dir)
	if window_target == nil {
		if focus_change_dir < 0 {
			window_target = c.FindWindowNavigable(len(c.Windows)-1, i_current, focus_change_dir)
		} else {
			window_target = c.FindWindowNavigable(0, i_current, focus_change_dir)
		}
	}
	c.NavWindowingTarget = window_target
	c.NavWindowingToggleLayer = false
}

// Equivalent of IsKeyDown() for NavInputs[]
func (c *Context) IsNavInputDown(n NavInput) bool {
	return c.IO.NavInputs[n] > 0.0
}

func (n *NavMoveResult) Clear() {
	n.ID = 0
	n.ParentID = 0
	n.Window = nil
	n.DistBox = math.MaxFloat32
	n.DistCenter = math.MaxFloat32
	n.DistAxial = math.MaxFloat32
	n.RectRel = f64.Rectangle{}
}

func (c *Context) NavUpdate() {
}