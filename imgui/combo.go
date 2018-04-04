package imgui

func (c *Context) EndCombo() {
	style := &c.Style
	if style.FramePadding.X != style.WindowPadding.X {
		c.UnindentEx(style.FramePadding.X - style.WindowPadding.X)
	}
	c.EndPopup()
}
