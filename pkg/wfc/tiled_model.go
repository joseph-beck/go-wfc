package wfc

import (
	"image"
	"image/color"
)

type TiledModel struct {
	*BaseModel
	TileSize   int
	Tiles      []TilePattern
	Propagator [][][]bool
}

type Tile struct {
	Name     string
	Sym      string
	Weight   float64
	Variants []image.Image
}

type Neighbour struct {
	Left     string
	LeftNum  int
	Right    string
	RightNum int
}

type TilePattern []color.Color

type Inversion func(int) int

func NewTiledModel(data TiledData, width int, height int, periodic bool) *TiledModel {

	// Initialize m
	m := &TiledModel{
		BaseModel: &BaseModel{
			Fmx:        width,
			Fmy:        height,
			Periodic:   periodic,
			Stationary: make([]float64, 0),
		},
		TileSize: data.TileSize,
		Tiles:    make([]TilePattern, 0),
	}

	first := make(map[string]int)
	action := make([][]int, 0)

	tile := func(transformer func(x int, y int) color.Color) TilePattern {
		result := make(TilePattern, m.TileSize*m.TileSize)
		for y := 0; y < m.TileSize; y++ {
			for x := 0; x < m.TileSize; x++ {
				result[x+y*m.TileSize] = transformer(x, y)
			}
		}
		return result
	}

	rotate := func(p TilePattern) TilePattern {
		return tile(func(x, y int) color.Color {
			return p[m.TileSize-1-y+x*m.TileSize]
		})
	}

	for i := 0; i < len(data.Tiles); i++ {
		current := data.Tiles[i]
		var (
			cardinality int
			inv1        Inversion
			inv2        Inversion
		)

		switch current.Sym {
		case "L":
			cardinality = 4
			inv1 = func(i int) int {
				return (i + 1) % 4
			}
			inv2 = func(i int) int {
				if i%2 == 0 {
					return i + 1
				}
				return i - 1
			}
		case "T":
			cardinality = 4
			inv1 = func(i int) int {
				return (i + 1) % 4
			}
			inv2 = func(i int) int {
				if i%2 == 0 {
					return i
				}
				return 4 - i
			}
		case "I":
			cardinality = 2
			inv1 = func(i int) int {
				return 1 - i
			}
			inv2 = func(i int) int {
				return i
			}
		case "\\":
			cardinality = 2
			inv1 = func(i int) int {
				return 1 - i
			}
			inv2 = func(i int) int {
				return 1 - i
			}
		case "X":
			cardinality = 1
			inv1 = func(i int) int {
				return i
			}
			inv2 = func(i int) int {
				return i
			}
		default:
			cardinality = 1
			inv1 = func(i int) int {
				return i
			}
			inv2 = func(i int) int {
				return i
			}
		}

		m.T = len(action)
		first[current.Name] = m.T

		for t := 0; t < cardinality; t++ {
			action = append(action, []int{
				m.T + t,
				m.T + inv1(t),
				m.T + inv1(inv1(t)),
				m.T + inv1(inv1(inv1(t))),
				m.T + inv2(t),
				m.T + inv2(inv1(t)),
				m.T + inv2(inv1(inv1(t))),
				m.T + inv2(inv1(inv1(inv1(t)))),
			})
		}

		if data.Unique {
			for t := 0; t < cardinality; t++ {
				img := current.Variants[t]
				m.Tiles = append(m.Tiles, tile(func(x, y int) color.Color {
					return img.At(x, y)
				}))
			}
		} else {
			img := current.Variants[0]
			m.Tiles = append(m.Tiles, tile(func(x, y int) color.Color {
				return img.At(x, y)
			}))

			for t := 1; t < cardinality; t++ {
				m.Tiles = append(m.Tiles, rotate(m.Tiles[m.T+t-1]))
			}
		}

		for t := 0; t < cardinality; t++ {
			m.Stationary = append(m.Stationary, current.Weight)
		}
	}

	m.T = len(action)
	m.Propagator = make([][][]bool, 4)

	for i := 0; i < 4; i++ {
		m.Propagator[i] = make([][]bool, m.T)
		for t := 0; t < m.T; t++ {
			m.Propagator[i][t] = make([]bool, m.T)
			for t2 := 0; t2 < m.T; t2++ {
				m.Propagator[i][t][t2] = false
			}
		}
	}

	m.Wave = make([][][]bool, m.Fmx)
	m.Changes = make([][]bool, m.Fmx)

	for x := 0; x < m.Fmx; x++ {
		m.Wave[x] = make([][]bool, m.Fmy)
		m.Changes[x] = make([]bool, m.Fmy)
		for y := 0; y < m.Fmy; y++ {
			m.Wave[x][y] = make([]bool, m.T)
		}
	}

	for i := 0; i < len(data.Neighbors); i++ {
		neighbor := data.Neighbors[i]

		l := action[first[neighbor.Left]][neighbor.LeftNum]
		d := action[l][1]
		r := action[first[neighbor.Right]][neighbor.RightNum]
		u := action[r][1]

		m.Propagator[0][r][l] = true
		m.Propagator[0][action[r][6]][action[l][6]] = true
		m.Propagator[0][action[l][4]][action[r][4]] = true
		m.Propagator[0][action[l][2]][action[r][2]] = true

		m.Propagator[1][u][d] = true
		m.Propagator[1][action[d][6]][action[u][6]] = true
		m.Propagator[1][action[u][4]][action[d][4]] = true
		m.Propagator[1][action[d][2]][action[u][2]] = true
	}

	for t := 0; t < m.T; t++ {
		for t2 := 0; t2 < m.T; t2++ {
			m.Propagator[2][t][t2] = m.Propagator[0][t2][t]
			m.Propagator[3][t][t2] = m.Propagator[1][t2][t]
		}
	}

	return m
}

func (model *TiledModel) OnBoundary(x int, y int) bool {
	return false
}

func (model *TiledModel) Propagate() bool {
	change := false

	for x2 := 0; x2 < model.Fmx; x2++ {
		for y2 := 0; y2 < model.Fmy; y2++ {
			for d := 0; d < 4; d++ {
				x1 := x2
				y1 := y2

				if d == 0 {
					if x2 == 0 {
						if !model.Periodic {
							continue
						} else {
							x1 = model.Fmx - 1
						}
					} else {
						x1 = x2 - 1
					}
				} else if d == 1 {
					if y2 == model.Fmy-1 {
						if !model.Periodic {
							continue
						} else {
							y1 = 0
						}
					} else {
						y1 = y2 + 1
					}
				} else if d == 2 {
					if x2 == model.Fmx-1 {
						if !model.Periodic {
							continue
						} else {
							x1 = 0
						}
					} else {
						x1 = x2 + 1
					}
				} else {
					if y2 == 0 {
						if !model.Periodic {
							continue
						} else {
							y1 = model.Fmy - 1
						}
					} else {
						y1 = y2 - 1
					}
				}

				if !model.Changes[x1][y1] {
					continue
				}

				for t2 := 0; t2 < model.T; t2++ {
					if model.Wave[x2][y2][t2] {
						b := false

						for t1 := 0; t1 < model.T && !b; t1++ {
							if model.Wave[x1][y1][t1] {
								b = model.Propagator[d][t2][t1]
							}
						}

						if !b {
							model.Wave[x2][y2][t2] = false
							model.Changes[x2][y2] = true
							change = true
						}
					}
				}
			}
		}
	}

	return change
}

func (m *TiledModel) Clear() {
	m.ClearBase(m)
}

func (model *TiledModel) RenderCompleteImage() image.Image {
	output := make([][]color.Color, model.Fmx*model.TileSize)
	for i := range output {
		output[i] = make([]color.Color, model.Fmy*model.TileSize)
	}

	for y := 0; y < model.Fmy; y++ {
		for x := 0; x < model.Fmx; x++ {
			for yt := 0; yt < model.TileSize; yt++ {
				for xt := 0; xt < model.TileSize; xt++ {
					for t := 0; t < model.T; t++ {
						if model.Wave[x][y][t] {
							output[x*model.TileSize+xt][y*model.TileSize+yt] = model.Tiles[t][yt*model.TileSize+xt]
							break
						}
					}
				}
			}
		}
	}
	return GeneratedImage{output}
}

func (model *TiledModel) RenderIncompleteImage() image.Image {
	output := make([][]color.Color, model.Fmx*model.TileSize)
	for i := range output {
		output[i] = make([]color.Color, model.Fmy*model.TileSize)
	}

	for y := 0; y < model.Fmy; y++ {
		for x := 0; x < model.Fmx; x++ {
			amount := 0
			sum := 0.0
			for t := 0; t < len(model.Wave[x][y]); t++ {
				if model.Wave[x][y][t] {
					amount += 1
					sum += model.Stationary[t]
				}
			}
			for yt := 0; yt < model.TileSize; yt++ {
				for xt := 0; xt < model.TileSize; xt++ {
					if amount == model.T {
						output[x*model.TileSize+xt][y*model.TileSize+yt] = color.RGBA{127, 127, 127, 255}
					} else {
						sR, sG, sB, sA := 0.0, 0.0, 0.0, 0.0
						for t := 0; t < model.T; t++ {
							if model.Wave[x][y][t] {
								r, g, b, a := model.Tiles[t][yt*model.TileSize+xt].RGBA()
								sR += float64(r) * model.Stationary[t]
								sG += float64(g) * model.Stationary[t]
								sB += float64(b) * model.Stationary[t]
								sA += float64(a) * model.Stationary[t]
							}
						}
						uR := uint8(int(sR/sum) >> 8)
						uG := uint8(int(sG/sum) >> 8)
						uB := uint8(int(sB/sum) >> 8)
						uA := uint8(int(sA/sum) >> 8)
						output[x*model.TileSize+xt][y*model.TileSize+yt] = color.RGBA{uR, uG, uB, uA}
					}
				}
			}
		}
	}
	return GeneratedImage{output}
}

func (model *TiledModel) Render() image.Image {
	if model.IsGenSuccess() {
		return model.RenderCompleteImage()
	} else {
		return model.RenderIncompleteImage()
	}
}

func (m *TiledModel) Iterate(iterations int) (image.Image, bool, bool) {
	finished := m.BaseModel.Iterate(m, iterations)
	return m.Render(), finished, m.IsGenSuccess()
}

func (m *TiledModel) Generate() (image.Image, bool) {
	m.BaseModel.Generate(m)
	return m.Render(), m.IsGenSuccess()
}
