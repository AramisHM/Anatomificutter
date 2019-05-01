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

	"github.com/nfnt/resize"
)

func printUsage() {
	fmt.Printf(`

	ANATOMIFICUTTER [v0.1.0]                             

USAGE

param 1 - Input directory

param 2 - Output directory
`)
	os.Exit(0)
}

// GetImageLen - Get image length from a file
func GetImageLen(filePath os.FileInfo) int {
	tmpFile, err := os.Open(filePath.Name())
	defer tmpFile.Close()
	if err != nil {
		fmt.Println("Image file couldn't be opened to get image length.")
		os.Exit(1)
	}
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

// GenerateCorCutOnZ
func GenerateCorCutOnZ(zLevel int, targetDir string, theFiles []os.FileInfo) {
	width := GetImageLen(theFiles[0])
	height := len(theFiles)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for y, f := range theFiles {
		fullPathFile := f.Name()
		fileExt := filepath.Ext(fullPathFile)
		if fileExt != ".png" {
			continue
		}

		imgfile, err := os.Open(fullPathFile)
		if err != nil {
			fmt.Println("Image file couldn't be opened.")
			imgfile.Close()
			continue
		}
		imgFromFile, _, err := image.Decode(imgfile)

		// Set color for each pixel.
		for x := 0; x < width; x++ {
			color := imgFromFile.At(x, zLevel)
			img.Set(x, y, color)
		}
		imgfile.Close()

	}
	//resize
	resizedImage := resize.Resize((uint)(width), (uint)(height/2), img, resize.Lanczos3)

	// Encode as PNG.

	f, _ := os.Create(targetDir + strconv.Itoa(zLevel) + ".png")
	png.Encode(f, resizedImage)
}

func main() {

	// Get arguments
	argsWithoutProg := os.Args[1:]

	// Get the parameters
	if len(argsWithoutProg) < 2 {
		printUsage()
	}

	counter, theFiles := GetFilesFromExtension("png", argsWithoutProg[0])

	path := "./out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if counter > 0 {
		for i := 0; i < counter; i++ {
			GenerateCorCutOnZ(i, "./out/", theFiles)
		}
	}
}
