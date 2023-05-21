package wfc

import (
	"image"
	"testing"
	"wfc/pkg/utils"
)

func overlappingTest(t *testing.T, filename, snapshotFilename string, iterations int) {
	periodicInput := true
	periodicOutput := true
	hasGround := true
	width := 48
	height := 48
	n := 3
	symetry := 2
	seed := int64(42)

	inputImg, err := utils.LoadImage("../../internal/input/" + filename)
	if err != nil {
		panic(err)
	}

	var outputImg image.Image
	success, finished := false, false
	model := NewOverlappingModel(inputImg, n, width, height, periodicInput, periodicOutput, symetry, hasGround)
	model.SetSeed(seed)

	if iterations == -1 {
		outputImg, success = model.Generate()
		if !success {
			t.Log("Failed to generate image on the first try.")
			t.FailNow()
		}
	} else {
		outputImg, finished, _ = model.Iterate(iterations)
		if finished {
			t.Log("Test for incomplete state actually finished.")
			t.FailNow()
		}
	}

	snapshotImg, err := utils.LoadImage("../../internal/snapshots/" + snapshotFilename)
	if err != nil {
		panic(err)
	}
	areEqual := utils.CompareImages(outputImg, snapshotImg)
	if !areEqual {
		t.Log("Output image is not the same as the snapshot image.")
		t.FailNow()
	}
}

func TestOverlappingGenerationCompletes(t *testing.T) {
	overlappingTest(t, "flowers.png", "flowers.png", -1)
}

func TestOverlappingIterationIncomplete(t *testing.T) {
	overlappingTest(t, "flowers.png", "flowers_incomplete.png", 5)
}
