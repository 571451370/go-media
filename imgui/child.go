package imgui

import "github.com/qeedquan/go-media/math/f64"

func (c *Context) BeginChild(str_id string, size_arg f64.Vec2, border bool, extra_flags WindowFlags) bool {
	window := c.GetCurrentWindow()
	return c.BeginChildEx(str_id, window.GetID(str_id), size_arg, border, extra_flags)
}

func (c *Context) BeginChildEx(name string, id ID, size_arg f64.Vec2, border bool, extra_flags WindowFlags) bool {
	return false
}

func (c *Context) EndChild() {
}

func (c *Context) BeginChildFrame(id ID, size f64.Vec2, extra_flags WindowFlags) bool {
	style := &c.Style
	c.PushStyleColor(ColChildBg, style.Colors[ColFrameBg].ToRGBA())
	c.PushStyleVar(StyleVarChildRounding, style.FrameRounding)
	c.PushStyleVar(StyleVarChildBorderSize, style.FrameBorderSize)
	c.PushStyleVar(StyleVarWindowPadding, style.FramePadding)
	return c.BeginChildEx("", id, size, true, WindowFlagsNoMove|WindowFlagsAlwaysUseWindowPadding|extra_flags)
}

func (c *Context) EndChildFrame() {
	c.EndChild()
	c.PopStyleVar(0)
	c.PopStyleColor()
}
