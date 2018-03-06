package imgui

import (
	"github.com/qeedquan/go-media/math/f64"
)

type IO struct {
	DisplaySize f64.Vec2
	DeltaTime   float64
	ConfigFlags ConfigFlags

	MetricsRenderVertices int
	MetricsRenderIndices  int
	MetricsActiveWindows  int

	KeyRepeatDelay float64
	KeyRepeatRate  float64

	FontDefault     *Font
	Fonts           *FontAtlas
	FontGlobalScale float64

	//------------------------------------------------------------------
	// Input - Fill before calling NewFrame()
	//------------------------------------------------------------------

	MousePos        f64.Vec2               // Mouse position, in pixels. Set to ImVec2(-FLT_MAX,-FLT_MAX) if mouse is unavailable (on another screen, etc.)
	MouseDown       [5]bool                // Mouse buttons: left, right, middle + extras. ImGui itself mostly only uses left button (BeginPopupContext** are using right button). Others buttons allows us to track if the mouse is being used by your application + available to user as a convenience via IsMouse** API.
	MouseWheel      float64                // Mouse wheel: 1 unit scrolls about 5 lines text.
	MouseWheelH     float64                // Mouse wheel (Horizontal). Most users don't have a mouse with an horizontal wheel, may not be filled by all back-ends.
	MouseDrawCursor bool                   // Request ImGui to draw a mouse cursor for you (if you are on a platform without a mouse cursor).
	KeyCtrl         bool                   // Keyboard modifier pressed: Control
	KeyShift        bool                   // Keyboard modifier pressed: Shift
	KeyAlt          bool                   // Keyboard modifier pressed: Alt
	KeySuper        bool                   // Keyboard modifier pressed: Cmd/Super/Windows
	KeysDown        [512]bool              // Keyboard keys that are pressed (ideally left in the "native" order your engine has access to keyboard keys, so you can use your own defines/enums for keys).
	InputCharacters [16 + 1]rune           // List of characters input (translated by user from keypress+keyboard state). Fill using AddInputCharacter() helper.
	NavInputs       [NavInputCount]float64 // Gamepad inputs (keyboard keys will be auto-mapped and be written here by ImGui::NewFrame)
}

type ConfigFlags uint
