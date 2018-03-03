package imgui

import (
	"hash/fnv"

	"github.com/qeedquan/go-media/math/f64"
)

type ID uint

type Context struct {
	IO    IO
	Style Style

	DrawDataBuilder DrawDataBuilder

	FontSize float64

	Windows       []*Window
	CurrentWindow *Window
	MovingWindow  *Window

	ActiveIdWindow           *Window
	ActiveId                 ID
	ActiveIdTimer            float64
	ActiveIdIsAlive          bool
	ActiveIdIsJustActivated  bool
	ActiveIdAllowOverlap     bool
	ActiveIdAllowNavDirFlags int
	ActiveIdClickOffset      f64.Vec2
	ActiveIdSource           InputSource

	NavWindow                  *Window
	NavId                      ID
	NavActivateId              ID
	NavActivateDownId          ID
	NavActivatePressedId       ID
	NavInputId                 ID
	NavJustTabbedId            ID
	NavNextActivateId          ID
	NavJustMovedToId           ID
	NavScoringRectScreen       f64.Rectangle
	NavScoringCount            int
	NavWindowingTarget         *Window
	NavWindowingHighlightTimer float64
	NavWindowingHighlightAlpha float64
	NavWindowingToggleLayer    bool
	NavWindowingInputSource    InputSource
	NavLayer                   int
	NavIdTabCounter            int
	NavIdIsAlive               bool
	NavMousePosDirty           bool
	NavDisableHighlight        bool
	NavDisableMouseHover       bool
	NavAnyRequest              bool
	NavInitRequest             bool
	NavInitRequestFromMove     bool
	NavInitResultId            ID
	NavInitResultRectRel       f64.Rectangle
	NavMoveFromClampedRefRect  bool
	NavMoveRequest             bool
	NavMoveRequestForward      NavForward
	NavMoveDir, NavMoveDirLast Dir
	NavMoveResultLocal         NavMoveResult
	NavMoveResultOther         NavMoveResult

	HoveredRootWindow     *Window
	HoveredWindow         *Window
	HoveredId             ID
	HoveredIdAllowOverlap bool
	HoveredIdTimer        float64

	Time               float64
	FrameCount         int
	FrameCountRendered int
	FrameCountEnded    int
	WindowsActiveCount int
}

func (c *Context) NewFrame() {
	c.FrameCount++
	c.WindowsActiveCount = 0
}

func (c *Context) EndFrame() {
	c.FrameCountEnded = c.FrameCount
}

func (c *Context) Render() {
	if c.FrameCountEnded != c.FrameCount {
		c.EndFrame()
	}
	c.FrameCountRendered = c.FrameCount

	// skip render altogether if alpha is 0
	if c.Style.Alpha <= 0 {
		return
	}

	io := &c.IO
	io.MetricsRenderVertices = 0
	io.MetricsRenderIndices = 0
	io.MetricsActiveWindows = 0

	// gather windows to render
	var wf *Window
	for _, w := range c.Windows {
		if w.Active && w.HiddenFrames <= 0 && w.Flags&WindowFlagsChildWindow == 0 && w != wf {
			c.AddWindowToDrawDataSelectLayer(w)
		}
	}
}

func (c *Context) AddWindowToDrawDataSelectLayer(window *Window) {
	io := c.GetIO()
	io.MetricsActiveWindows++
	if window.Flags&WindowFlagsTooltip != 0 {
		c.AddWindowToDrawData(&c.DrawDataBuilder.Layers[1], window)
	} else {
		c.AddWindowToDrawData(&c.DrawDataBuilder.Layers[0], window)
	}
}

func (c *Context) AddWindowToDrawData(outRenderList *[]*DrawList, window *Window) {
	dc := &window.DC
	c.AddDrawListToDrawData(outRenderList, window.DrawList)
	for _, child := range dc.ChildWindows {
		// clipped children may have been marked as not active
		if child.Active && child.HiddenFrames <= 0 {
			c.AddWindowToDrawData(outRenderList, child)
		}
	}
}

func (c *Context) AddDrawListToDrawData(outRenderList *[]*DrawList, drawList *DrawList) {
	n := len(drawList.CmdBuffer)
	if n == 0 {
		return
	}

	// remove trailing command if unused
	lastCmd := &drawList.CmdBuffer[n-1]
	if lastCmd.ElemCount == 0 && lastCmd.UserCallback == nil {
		drawList.CmdBuffer = drawList.CmdBuffer[:n-1]
		if len(drawList.CmdBuffer) == 0 {
			return
		}
	}

	*outRenderList = append(*outRenderList, drawList)
}

func (c *Context) GetIO() *IO {
	return &c.IO
}

func (c *Context) GetCurrentWindow() *Window {
	return c.CurrentWindow
}

func (c *Context) GetStyle() *Style {
	return &c.Style
}

func (c *Context) GetActiveID() ID {
	return c.ActiveId
}

func (c *Context) SetActiveID(id ID, window *Window) {
	c.ActiveIdIsJustActivated = c.ActiveId != id
	if c.ActiveIdIsJustActivated {
		c.ActiveIdTimer = 0
	}
	c.ActiveId = id
	c.ActiveIdAllowNavDirFlags = 0
	c.ActiveIdAllowOverlap = false
	c.ActiveIdWindow = window
	if id != 0 {
		c.ActiveIdIsAlive = true
		if c.NavActivateId == id || c.NavInputId == id || c.NavJustTabbedId == id || c.NavJustMovedToId == id {
			c.ActiveIdSource = InputSourceNav
		} else {
			c.ActiveIdSource = InputSourceMouse
		}
	}
}

func (c *Context) ClearActiveID() {
	c.SetActiveID(0, nil)
}

func (c *Context) SetHoveredID(id ID) {
	c.HoveredId = id
	c.HoveredIdAllowOverlap = false
}

func (c *Context) KeepAliveID(id ID) {
	if c.ActiveId == id {
		c.ActiveIdIsAlive = true
	}
}

func (c *Context) Hash(b []byte) uint32 {
	h := fnv.New32()
	h.Write(b)
	return h.Sum32()
}
