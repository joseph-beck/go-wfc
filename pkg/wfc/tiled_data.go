package wfc

import (
	"encoding/json"
	"image"
	"os"
	"strconv"
	"wfc/pkg/utils"
)

type RawData struct {
	Path       string         `json:"path"`      // Path to tiles
	Unique     bool           `json:"unique"`    // Default to false
	TileSize   int            `json:"tileSize"`  // Default to 16
	Tiles      []RawTile      `json:"tiles"`     //
	Neighbours []RawNeighbour `json:"neighbors"` //
}

type RawTile struct {
	Name     string  `json:"name"`     // Name used to identify the tile
	Symmetry string  `json:"symmetry"` // Default to ""
	Weight   float64 `json:"weight"`   // Default to 1
}

// Information on which tiles can be neighbors
type RawNeighbour struct {
	Left     string `json:"left"`     // Mathces Tile.Name
	LeftNum  int    `json:"leftNum"`  // Default to 0
	Right    string `json:"right"`    // Mathces Tile.Name
	RightNum int    `json:"rightNum"` // Default to 0
}

type TiledData struct {
	Unique    bool
	TileSize  int
	Tiles     []Tile
	Neighbors []Neighbour
}

func MakeTiledData(file string) TiledData {
	dataFile, err := os.ReadFile("../../internal/input/" + file)
	if err != nil {
		panic(err)
	}

	// Parse rd file
	var rd RawData
	if err := json.Unmarshal(dataFile, &rd); err != nil {
		panic(err)
	}

	// Marshal into data settings struct
	tiles := make([]Tile, len(rd.Tiles))
	for i, rt := range rd.Tiles {
		imgs := make([]image.Image, 0)
		if rd.Unique {
			i := 1
			for {
				if img, err := utils.LoadImage("../../internal/input/" + rd.Path + rt.Name + " " + strconv.Itoa(i) + ".png"); err == nil {
					imgs = append(imgs, img)
				} else {
					break
				}
				i++
			}
		} else {
			img, err := utils.LoadImage("../../internal/input/" + rd.Path + rt.Name + ".png")
			if err != nil {
				panic(err)
			}
			imgs = append(imgs, img)
		}

		weight := rt.Weight
		if weight == 0 {
			weight = 1
		}

		tiles[i] = Tile{
			Name:     rt.Name,
			Sym:      rt.Symmetry,
			Weight:   weight,
			Variants: imgs,
		}
	}

	neighbours := make([]Neighbour, len(rd.Neighbours))
	for i, rn := range rd.Neighbours {
		neighbours[i] = Neighbour(rn)
	}

	return TiledData{
		Unique:    rd.Unique,
		TileSize:  rd.TileSize,
		Tiles:     tiles,
		Neighbors: neighbours,
	}
}
