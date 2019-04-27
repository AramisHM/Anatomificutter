package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func printUsage() {
	fmt.Printf(`

	ANATOMIFICUTTER [v0.1.0]                             

USAGE

param 1 - Project Name

param 2 - Input directory

param 3 - Output directory
\n\n`)
	os.Exit(0)
}

func getFiles() {

}

func main() {

	// get arguments
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 2 {
		printUsage()
	}

	files, err := ioutil.ReadDir(argsWithoutProg[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()
		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		outputDir := "./"
		if len(argsWithoutProg) > 1 {
			outputDir = argsWithoutProg[2]
		}

		file, err := os.Create(outputDir + "/" + name + ".xml") // Truncates if file already exists, be careful!
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
			return
		}

		file.WriteString(`aio`)
		file.Close()
	}

	width := 300
	height := 300

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
