package wfc

import (
	"image"
	"image/color"
	"math"
)

type OverlappingModel struct {
	*BaseModel               // Base model
	N          int           // Pattern size
	Colors     []color.Color // Colours array
	Ground     int           // Pattern Id
	Patterns   []Pattern     // Unique pattern Ids from input
	Propagator [][][][]int   // Table of which patterns (t2) mathch a given pattern (t1) at offset (dx, dy) [t1][dx][dy][t2]
	Fmxmn      int           // Width - n
	Fmymn      int           // Height - n
}

type Pattern []int

func NewOverlappingModel(img image.Image, n, width, height int, periodicInput, periodic bool, symmetry int, ground bool) *OverlappingModel {
	m := &OverlappingModel{
		BaseModel: &BaseModel{
			Fmx:      width,
			Fmy:      height,
			Periodic: periodic,
		},
		N:      n,
		Ground: -1,
	}

	bounds := img.Bounds()
	dataWidth := bounds.Max.X
	dataHeight := bounds.Max.Y

	sample := make([][]int, dataWidth)
	for i := range sample {
		sample[i] = make([]int, dataHeight)
	}

	m.Colors = make([]color.Color, 0)
	colorMap := make(map[color.Color]int)

	for y := 0; y < dataHeight; y++ {
		for x := 0; x < dataWidth; x++ {
			color := img.At(x, y)
			if _, ok := colorMap[color]; !ok {
				colorMap[color] = len(m.Colors)
				m.Colors = append(m.Colors, color)
			}
			sample[x][y] = colorMap[color]
		}
	}

	c := len(m.Colors)
	w := int(math.Pow(float64(c), float64(n*n)))

	getPattern := func(transformer func(x, y int) int) Pattern {
		result := make(Pattern, n*n)
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				result[x+y*n] = transformer(x, y)
			}
		}
		return result
	}

	patternFromSample := func(x, y int) Pattern {
		return getPattern(func(dx, dy int) int {
			return sample[(x+dx)%dataWidth][(y+dy)%dataHeight]
		})
	}

	rotate := func(p Pattern) Pattern {
		return getPattern(func(x, y int) int {
			return p[n-1-y+x*n]
		})
	}

	reflect := func(p Pattern) Pattern {
		return getPattern(func(x, y int) int {
			return p[n-1-x+y*n]
		})
	}

	indexFromPattern := func(p Pattern) int {
		result := 0
		power := 1
		for i := 0; i < len(p); i++ {
			result += p[len(p)-1-i] * power
			power *= c
		}
		return result
	}

	patternFromIndex := func(ind int) Pattern {
		residue := ind
		power := w
		result := make(Pattern, n*n)
		for i := 0; i < len(result); i++ {
			power /= c
			count := 0
			for residue >= power {
				residue -= power
				count++
			}
			result[i] = count
		}
		return result
	}

	weights := make(map[int]int)
	weightsKeys := make([]int, 0)

	var (
		horizontalBound int
		verticalBound   int
	)

	if periodicInput {
		horizontalBound = dataWidth
		verticalBound = dataHeight
	} else {
		horizontalBound = dataWidth - n + 1
		verticalBound = dataHeight - n + 1
	}
	for y := 0; y < verticalBound; y++ {
		for x := 0; x < horizontalBound; x++ {
			ps := make([]Pattern, 8, 8)
			ps[0] = patternFromSample(x, y)
			ps[1] = reflect(ps[0])
			ps[2] = rotate(ps[0])
			ps[3] = reflect(ps[2])
			ps[4] = rotate(ps[2])
			ps[5] = reflect(ps[4])
			ps[6] = rotate(ps[4])
			ps[7] = reflect(ps[6])
			for k := 0; k < symmetry; k++ {
				ind := indexFromPattern(ps[k])
				if _, ok := weights[ind]; ok {
					weights[ind]++
				} else {
					weightsKeys = append(weightsKeys, ind)
					weights[ind] = 1
				}
				if ground && y == verticalBound-1 && x == 0 && k == 0 {
					// Set groung pattern
					m.Ground = len(weightsKeys) - 1
				}
			}
		}
	}

	m.T = len(weightsKeys)

	m.Patterns = make([]Pattern, m.T)
	m.Stationary = make([]float64, m.T)
	m.Propagator = make([][][][]int, m.T)

	for i, wk := range weightsKeys {
		m.Patterns[i] = patternFromIndex(wk)
		m.Stationary[i] = float64(weights[wk])
	}

	m.Wave = make([][][]bool, m.Fmx)
	m.Changes = make([][]bool, m.Fmx)

	for x := 0; x < m.Fmx; x++ {
		m.Wave[x] = make([][]bool, m.Fmy)
		m.Changes[x] = make([]bool, m.Fmy)

		for y := 0; y < m.Fmy; y++ {
			m.Wave[x][y] = make([]bool, m.T)
			m.Changes[x][y] = false

			for t := 0; t < m.T; t++ {
				m.Wave[x][y][t] = true
			}
		}
	}

	agrees := func(p1, p2 Pattern, dx, dy int) bool {
		var xmin, xmax, ymin, ymax int

		if dx < 0 {
			xmin = 0
			xmax = dx + n
		} else {
			xmin = dx
			xmax = n
		}

		if dy < 0 {
			ymin = 0
			ymax = dy + n
		} else {
			ymin = dy
			ymax = n
		}

		for y := ymin; y < ymax; y++ {
			for x := xmin; x < xmax; x++ {
				if p1[x+n*y] != p2[x-dx+n*(y-dy)] {
					return false
				}
			}
		}

		return true
	}

	for t := 0; t < m.T; t++ {
		m.Propagator[t] = make([][][]int, 2*n-1)
		for x := 0; x < 2*n-1; x++ {
			m.Propagator[t][x] = make([][]int, 2*n-1)
			for y := 0; y < 2*n-1; y++ {
				list := make([]int, 0)

				for t2 := 0; t2 < m.T; t2++ {
					if agrees(m.Patterns[t], m.Patterns[t2], x-n+1, y-n+1) {
						list = append(list, t2)
					}
				}

				m.Propagator[t][x][y] = make([]int, len(list))

				copy(m.Propagator[t][x][y], list)
			}
		}
	}

	m.Fmxmn = m.Fmx - m.N
	m.Fmymn = m.Fmy - m.N

	return m
}

func (m *OverlappingModel) OnBoundary(x, y int) bool {
	return !m.Periodic && (x > m.Fmxmn || y > m.Fmymn)
}

func (m *OverlappingModel) Propagate() bool {
	change := false
	startLoop := -m.N + 1
	endLoop := m.N

	for x := 0; x < m.Fmx; x++ {
		for y := 0; y < m.Fmy; y++ {
			if m.Changes[x][y] {
				m.Changes[x][y] = false

				for dx := startLoop; dx < endLoop; dx++ {
					for dy := startLoop; dy < endLoop; dy++ {
						sx := x + dx
						sy := y + dy

						if sx < 0 {
							sx += m.Fmx
						} else if sx >= m.Fmx {
							sx -= m.Fmx
						}

						if sy < 0 {
							sy += m.Fmy
						} else if sy >= m.Fmy {
							sy -= m.Fmy
						}

						if !m.Periodic && (sx > m.Fmx || sy > m.Fmy) {
							continue
						}

						allowed := m.Wave[sx][sy]

						for t := 0; t < m.T; t++ {
							if !allowed[t] {
								continue
							}

							b := false
							prop := m.Propagator[t][m.N-1-dx][m.N-1-dy]
							for i := 0; i < len(prop) && !b; i++ {
								b = m.Wave[x][y][prop[i]]
							}

							if !b {
								m.Changes[sx][sy] = true
								change = true
								allowed[t] = false
							}
						}
					}
				}
			}
		}
	}

	return change
}

func (m *OverlappingModel) Clear() {
	m.ClearBase(m)
	if m.Ground != -1 && m.T > 1 {
		for x := 0; x < m.Fmx; x++ {
			for t := 0; t < m.T; t++ {
				if t != m.Ground {
					m.Wave[x][m.Fmy-1][t] = false
				}
			}

			m.Changes[x][m.Fmy-1] = true

			for y := 0; y < m.Fmy-1; y++ {
				m.Wave[x][y][m.Ground] = false
				m.Changes[x][y] = true
			}
		}

		for m.Propagate() {
			// Empty loop
		}
	}
}

func (model *OverlappingModel) RenderCompleteImage() image.Image {
	output := make([][]color.Color, model.Fmx)
	for i := range output {
		output[i] = make([]color.Color, model.Fmy)
	}

	for y := 0; y < model.Fmy; y++ {
		for x := 0; x < model.Fmx; x++ {
			for t := 0; t < model.T; t++ {
				if model.Wave[x][y][t] {
					output[x][y] = model.Colors[model.Patterns[t][0]]
				}
			}
		}
	}

	return GeneratedImage{output}
}

func (m *OverlappingModel) RenderIncompleteImage() image.Image {
	output := make([][]color.Color, m.Fmx)
	for i := range output {
		output[i] = make([]color.Color, m.Fmy)
	}

	var (
		contributorNumber uint32
		sR                uint32
		sG                uint32
		sB                uint32
		sA                uint32
	)

	for y := 0; y < m.Fmy; y++ {
		for x := 0; x < m.Fmx; x++ {
			contributorNumber, sR, sG, sB, sA = 0, 0, 0, 0, 0

			for dy := 0; dy < m.N; dy++ {
				for dx := 0; dx < m.N; dx++ {
					sx := x - dx
					if sx < 0 {
						sx += m.Fmx
					}

					sy := y - dy
					if sy < 0 {
						sy += m.Fmy
					}

					if !m.Periodic && (sx > m.Fmxmn || sy > m.Fmymn) {
						continue
					}

					for t := 0; t < m.T; t++ {
						if m.Wave[sx][sy][t] {
							contributorNumber++
							r, g, b, a := m.Colors[m.Patterns[t][dx+dy*m.N]].RGBA()
							sR += r
							sG += g
							sB += b
							sA += a
						}
					}
				}
			}

			if contributorNumber == 0 {
				output[x][y] = color.RGBA{127, 127, 127, 255}
			} else {
				uR := uint8((sR / contributorNumber) >> 8)
				uG := uint8((sG / contributorNumber) >> 8)
				uB := uint8((sB / contributorNumber) >> 8)
				uA := uint8((sA / contributorNumber) >> 8)
				output[x][y] = color.RGBA{uR, uG, uB, uA}
			}
		}
	}
	return GeneratedImage{
		data: output,
	}
}

func (m *OverlappingModel) Render() image.Image {
	if m.GenSuccess {
		return m.RenderCompleteImage()
	} else {
		return m.RenderIncompleteImage()
	}
}

func (m *OverlappingModel) Iterate(iterations int) (image.Image, bool, bool) {
	finished := m.BaseModel.Iterate(m, iterations)
	return m.Render(), finished, m.GenSuccess
}

func (m *OverlappingModel) Generate() (image.Image, bool) {
	m.BaseModel.Generate(m)
	return m.Render(), m.GenSuccess
}
