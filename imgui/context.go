package imgui

import (
	"hash/fnv"

	"github.com/qeedquan/go-media/math/f64"
)

type ID uint

type Context struct {
	IO    IO
	Style Style

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

type DrawContext struct {
	CursorPos                 f64.Vec2
	CurrentLineTextBaseOffset float64
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
			c.AddWindowToDrawSelectLayer(w)
		}
	}
}

func (c *Context) AddWindowToDrawSelectLayer(w *Window) {
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

func (c *Context) SetActiveID(id ID, w *Window) {
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
