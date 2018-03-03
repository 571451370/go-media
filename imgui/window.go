package imgui

type WindowFlags uint

const (
	WindowFlagsNoInputs WindowFlags = 1 << iota
	WindowFlagsChildWindow
)

type Window struct {
	DC           DrawContext
	Active       bool
	HiddenFrames int
	Flags        WindowFlags
	SkipItems    bool
}
