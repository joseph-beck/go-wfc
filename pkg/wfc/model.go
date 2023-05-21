package wfc

import "image"

type Iterator interface {
	Iterate(iterations int) (image.Image, bool, bool)
}

type Generator interface {
	Generate() (image.Image, bool)
}

type Checker interface {
	OnBoundary(x int, y int) bool
}

type Propagater interface {
	Propagate() bool
}

type Clearer interface {
	Clear()
}

type Collapser interface {
	Iterator
	Generator
	Checker
	Propagater
	Clearer
}
