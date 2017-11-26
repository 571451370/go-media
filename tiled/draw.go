package tiled

import "image"

type Drawer struct {
	*Map
	Camera   image.Rectangle
	Viewport image.Rectangle
}

func (c *Drawer) Draw() {
	for i := len(c.Layers) - 1; i >= 0; i-- {
		c.DrawLayer(i)
	}
}

func (c *Drawer) DrawLayer(layer int) {
}