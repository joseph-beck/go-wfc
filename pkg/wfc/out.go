package wfc

import (
	"image"
	"image/color"
)

type GeneratedImage struct {
	data [][]color.Color
}

func (g GeneratedImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (g GeneratedImage) Bounds() image.Rectangle {
	if len(g.data) < 1 {
		return image.Rect(0, 0, 0, 0)
	}
	return image.Rect(0, 0, len(g.data[0]), len(g.data))
}

func (g GeneratedImage) At(x, y int) color.Color {
	return g.data[x][y]
}