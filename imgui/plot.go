package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) PlotLines(label string, values []float64) {
	c.PlotLinesEx(label, values, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0}, 4)
}

func (c *Context) PlotLinesEx(label string, values []float64, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2, stride int) {
}

func (c *Context) PlotLinesItem(label string, values_getter func(idx int) float64, values_count int) {
	c.PlotLinesItemEx(label, values_getter, values_count, 0, "", math.MaxFloat32, math.MaxFloat32, f64.Vec2{0, 0})
}

func (c *Context) PlotLinesItemEx(label string, values_getter func(idx int) float64, values_count, values_offset int, overlay_text string, scale_min, scale_max float64, graph_size f64.Vec2) {
}
