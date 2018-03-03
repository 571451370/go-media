package imgui

type SliderFlags uint

const (
	SliderFlagsVertical SliderFlags = 1 << iota
)

func (c *Context) SliderInt() {
}

func (c *Context) VSliderInt(label string, vals []int, vmin, vmax int, format string) {
}
