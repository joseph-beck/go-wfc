package wfc

import (

	// "fmt"
	"image"
	"testing"

	"wfc/pkg/utils"
)

func simpleTiledTest(t *testing.T, dataFileName, snapshotFileName string, iterations int) {
	// Set test parameters
	periodic := false
	width := 20
	height := 20
	seed := int64(42)
	data := MakeTiledData("../../internal/input/", dataFileName)

	var outputImg image.Image
	success, finished := false, false
	model := NewTiledModel(data, width, height, periodic)
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

	snapshotImg, err := utils.LoadImage("../../internal/snapshots/" + snapshotFileName)
	if err != nil {
		panic(err)
	}

	areEqual := utils.CompareImages(outputImg, snapshotImg)
	if !areEqual {
		t.Log("Output image is not the same as the snapshot image.")
		t.FailNow()
	}
}

func TestSimpleTiledGenerationCompletes(t *testing.T) {
	simpleTiledTest(t, "castle_data.json", "castle.png", -1)
}

func TestSimpleTiledIterationIncomplete(t *testing.T) {
	simpleTiledTest(t, "castle_data.json", "castle_incomplete.png", 5)
}
