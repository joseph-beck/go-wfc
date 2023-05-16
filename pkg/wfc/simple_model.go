package wfc

import (
	"image"
	"image/color"
)

type SimpleTiledModel struct {
	*BaseModel
	TileSize   int
	Tiles      []TilePattern
	Propagater [][][]bool
}

type SimpleTiledData struct {
	Unique    bool
	TileSize  int
	Tiles     []Tile
	Neighbors []Neighbor
}

type Tile struct {
	Name     string
	Sym      string
	Weight   float64
	Variants []image.Image
}

type Neighbor struct {
	Left     string
	LeftNum  int
	Right    string
	RightNum int
}

type TilePattern []color.Color
