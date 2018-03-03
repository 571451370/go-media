package imgui

import "github.com/qeedquan/go-media/math/f64"

type WindowFlags uint

const (
	WindowFlagsNoTitleBar                WindowFlags = 1 << iota // Disable title-bar
	WindowFlagsNoResize                                          // Disable user resizing with the lower-right grip
	WindowFlagsNoMove                                            // Disable user moving the window
	WindowFlagsNoScrollbar                                       // Disable scrollbars (window can still scroll with mouse or programatically)
	WindowFlagsNoScrollWithMouse                                 // Disable user vertically scrolling with mouse wheel. On child window mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	WindowFlagsNoCollapse                                        // Disable user collapsing window by double-clicking on it
	WindowFlagsAlwaysAutoResize                                  // Resize every window to its content every frame
	WindowFlagsNoSavedSettings                                   // Never load/save settings in .ini file
	WindowFlagsNoInputs                                          // Disable catching mouse or keyboard inputs hovering test with pass through.
	WindowFlagsMenuBar                                           // Has a menu-bar
	WindowFlagsHorizontalScrollbar                               // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing                                // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus                             // Disable bringing window to front when taking focus (e.g. clicking on it or programatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar                           // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar                         // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding                            // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows because more convenient)
	WindowFlagsResizeFromAnySide                                 // (WIP) Enable resize from any corners and borders. Your back-end needs to honor the different values of io.MouseCursor set by imgui.
	WindowFlagsNoNavInputs                                       // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                                        // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	WindowFlagsNoNav                     WindowFlags = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus

	// [Internal]
	WindowFlagsNavFlattened WindowFlags = 1 << iota // (WIP) Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	WindowFlagsChildWindow                          // Don't use! For internal use by BeginChild()
	WindowFlagsTooltip                              // Don't use! For internal use by BeginTooltip()
	WindowFlagsPopup                                // Don't use! For internal use by BeginPopup()
	WindowFlagsModal                                // Don't use! For internal use by BeginPopupModal()
	WindowFlagsChildMenu                            // Don't use! For internal use by BeginMenu()
)

type Window struct {
	Name                string
	Id                  ID
	Flags               WindowFlags
	PosFloat            f64.Vec2
	Pos                 f64.Vec2
	Size                f64.Vec2
	SizeFull            f64.Vec2
	SizeFullAtLastBegin f64.Vec2
	SizeContents        f64.Vec2
	DC                  DrawContext
	Active              bool
	WasActive           bool
	HiddenFrames        int
	SkipItems           bool
	DrawList            *DrawList
}

type DrawContext struct {
	CursorPos                 f64.Vec2
	CurrentLineTextBaseOffset float64
	ChildWindows              []*Window
}