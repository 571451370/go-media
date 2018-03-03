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
}

type ConfigFlags uint
