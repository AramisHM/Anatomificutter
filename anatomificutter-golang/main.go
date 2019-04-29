package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func printUsage() {
	fmt.Printf(`

	ANATOMIFICUTTER [v0.1.0]                             

USAGE

param 1 - Input directory

param 2 - Coronal cut depth
`)
	os.Exit(0)
}

// GetImageLen - Get image length from a file
func GetImageLen(filePath os.FileInfo) int {
	tmpFile, err := os.Open(filePath.Name())
	if err != nil {
		fmt.Println("Image file couldn't be opened.")
		os.Exit(1)
	}
	defer tmpFile.Close()
	imageToGetLenth, _, err := image.Decode(tmpFile)
	lineLen := imageToGetLenth.Bounds().Dx()

	return lineLen
}

// GetFilesFromExtension
func GetFilesFromExtension(extension string, path string) (int, []os.FileInfo) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// get number of png files
	var theFiles []os.FileInfo
	counter := 0
	for _, f := range files {
		fullPathFile := f.Name()
		fileExt := filepath.Ext(fullPathFile)

		// only .png files allowed
		if fileExt != "."+extension {
			continue
		} else {
			theFiles = append(theFiles, f)
			counter++
		}
	}

	return counter, theFiles
}

func main() {

	// Get arguments
	argsWithoutProg := os.Args[1:]

	// Get the parameters
	if len(argsWithoutProg) < 2 {
		printUsage()
	}

	cutDepth, _ := strconv.ParseInt(argsWithoutProg[1], 10, 64)

	counter, theFiles := GetFilesFromExtension("png", argsWithoutProg[0])

	if counter > 0 {
		// Get image length
		//lineHeight := 0
		lineLen := GetImageLen(theFiles[0])

		width := lineLen
		height := counter

		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}
		img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

		for y, f := range theFiles {
			fullPathFile := f.Name()
			fileExt := filepath.Ext(fullPathFile)
			fileNameOnly := strings.TrimSuffix(fullPathFile, filepath.Ext(fullPathFile))

			// only .png files allowed
			if fileExt != ".png" {
				continue
			}

			imgfile, err := os.Open(fullPathFile)
			if err != nil {
				fmt.Println("Image file couldn't be opened.")
				os.Exit(1)
			}
			defer imgfile.Close()

			imgFromFile, _, err := image.Decode(imgfile)

			// Set color for each pixel.
			for x := 0; x < width; x++ {
				color := imgFromFile.At(x, (int)(cutDepth))
				img.Set(x, y, color)
			}

			// Encode as PNG.
			f, _ := os.Create("out-" + fileNameOnly + ".png")
			png.Encode(f, img)
		}
	}
}
