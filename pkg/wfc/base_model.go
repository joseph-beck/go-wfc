package wfc

import "math"

type BaseModel struct {
	InitField  bool           // Generation initiliazed?
	RngSet     bool           // Random number generator set by user?
	GenSuccess bool           // Generation has run into a contradiction?
	Wave       [][][]bool     // All possible patterns (t) that could fit (x, y)
	Changes    [][]bool       // Chabges made in interation of propagation
	Stationary []float64      // Array of weights for patterns
	T          int            // Count of patterns
	Periodic   bool           // Tessellates?
	Fmx        int            // Width
	Fmy        int            // Height
	Rng        func() float64 // Random number generator supplied at gen time
}

func (b *BaseModel) Observe(sm AlgorithmApplier) bool {
	min := 1000.0
	argminx := -1
	argminy := -1
	distribution := make([]float64, b.T)

	for x := 0; x < b.Fmx; x++ {
		for y := 0; y < b.Fmy; y++ {
			if sm.OnBoundary(x, y) {
				continue
			}

			sum := 0.0

			for t := 0; t < b.T; t++ {
				if b.Wave[x][y][t] {
					distribution[t] = b.Stationary[t]
				} else {
					distribution[t] = 0.0
				}
				sum += distribution[t]
			}

			if sum == 0.0 {
				b.GenSuccess = false
				return true // finished, unsuccessful
			}

			for t := 0; t < b.T; t++ {
				distribution[t] /= sum
			}

			entropy := 0.0

			for i := 0; i < len(distribution); i++ {
				if distribution[i] > 0.0 {
					entropy += -distribution[i] * math.Log(distribution[i])
				}
			}

			noise := 0.000001 * b.Rng()

			if entropy > 0 && entropy+noise < min {
				min = entropy + noise
				argminx = x
				argminy = y
			}
		}
	}

	if argminx == -1 && argminy == -1 {
		b.GenSuccess = true
		return true
	}

	for t := 0; t < b.T; t++ {
		if b.Wave[argminx][argminy][t] {
			distribution[t] = b.Stationary[t]
		} else {
			distribution[t] = 0.0
		}
	}

	r := randomIndice(distribution, b.Rng())

	for t := 0; t < b.T; t++ {
		b.Wave[argminx][argminy][t] = (t == r)
	}

	b.Changes[argminx][argminy] = true

	return false
}
