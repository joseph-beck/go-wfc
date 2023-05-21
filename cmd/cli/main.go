package main

import (
	"image/jpeg"
	"log"
	"os"
	"wfc/pkg/utils"
	"wfc/pkg/wfc"
)

func main() {
	data := wfc.MakeTiledData("pipe_data.json")
	t := wfc.NewTiledModel(data, 20, 20, false)

	i, _ := t.Generate()
	// with Makefile remove ../../
	f, err := os.Create("../../internal/output/img.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	err = jpeg.Encode(f, i, nil)
	if err != nil {
		log.Fatalln(err)
	}

	img, err := utils.LoadImage("../../internal/input/maze.png")
	if err != nil {
		log.Fatalln(err)
	}
	o := wfc.NewOverlappingModel(img, 3, 16, 16, true, true, 2, false)

	i, _ = o.Generate()
	f, err = os.Create("../../internal/output/img2.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	err = jpeg.Encode(f, i, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
