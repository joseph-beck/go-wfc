package main

import (
	"image/jpeg"
	"log"
	"os"
	"wfc/pkg/wfc"
)

func main() {
	data := wfc.MakeTiledData("pipe_data.json")
	m := wfc.NewTiledModel(data, 20, 20, false)

	i, _ := m.Generate()
	// with Makefile remove ../../
	f, err := os.Create("../../internal/output/img.jpg")

	if err != nil {
		log.Fatalln(err)
	}
	err = jpeg.Encode(f, i, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
