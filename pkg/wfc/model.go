package wfc

import "image"

type Iterator interface {
	Iterator(iterations int) (image.Image, bool, bool)
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

type AlgorithmApplier interface {
	Iterator
	Generator
	Checker
	Propagater
	Clearer
}
