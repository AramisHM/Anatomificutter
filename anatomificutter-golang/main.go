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

param 1 - input path

param 2 - "sagital" or "coronal"
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

// GetImageHeight - Get image height from a file
func GetImageHeight(filePath os.FileInfo) int {
	tmpFile, err := os.Open(filePath.Name())
	defer tmpFile.Close()
	if err != nil {
		fmt.Println("Image file couldn't be opened to get image length.")
		os.Exit(1)
	}
	imageToGetLenth, _, err := image.Decode(tmpFile)
	lineHeight := imageToGetLenth.Bounds().Dy()

	return lineHeight
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

	f, _ := os.Create(targetDir + "/" + strconv.Itoa(zLevel) + ".png")
	png.Encode(f, resizedImage)
}

// GenerateCorCutOnZ
func GenerateSagCutOnX(xLevel int, targetDir string, theFiles []os.FileInfo) {
	width := GetImageHeight(theFiles[0])
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
			color := imgFromFile.At(xLevel, x)
			img.Set(x, y, color)
		}
		imgfile.Close()

	}
	//resize
	resizedImage := resize.Resize((uint)(width), (uint)(height/2), img, resize.Lanczos3)

	// Encode as PNG.

	f, _ := os.Create(targetDir + "/" + strconv.Itoa(xLevel) + ".png")
	png.Encode(f, resizedImage)
}

func DoCoronal(theFiles []os.FileInfo) {
	path := "./coronal-out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if len(theFiles) > 0 {

		numCuts := GetImageHeight(theFiles[0])
		for i := 0; i < numCuts; i++ {
			GenerateCorCutOnZ(i, path, theFiles)
		}
	}
}

func DoSagital(theFiles []os.FileInfo) {
	path := "./sagital-out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if len(theFiles) > 0 {

		numCuts := GetImageLen(theFiles[0])
		for i := 0; i < numCuts; i++ {
			GenerateCorCutOnZ(i, path, theFiles)
		}
	}
}
func main() {

	// Get arguments
	argsWithoutProg := os.Args[1:]

	// Get the parameters
	if len(argsWithoutProg) < 1 {
		printUsage()
	}

	counter, theFiles := GetFilesFromExtension("png", argsWithoutProg[0])

	if counter > 0 {
		processType := argsWithoutProg[1]

		switch processType {
		case "coronal":
			DoCoronal(theFiles)
		case "sagital":
			DoSagital(theFiles)
		}
	}

}
