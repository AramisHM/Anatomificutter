package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func printUsage() {
	fmt.Printf(`

	ANATOMIFICUTTER [v0.1.0]                             

USAGE

param 1 - input path

param 2 - "sagital" or "coronal"

param 3 - max RAM memory to use in Mega Bytes (MB)
`)
	os.Exit(0)
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
func CloneToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

// GenerateCorCutOnZ
func GenerateCorCutOnZ(zLevel int, startYIndex int, targetDir string, imgs []image.Image, finalHeight int) {
	width := imgs[0].Bounds().Dx()
	height := finalHeight

	targFileName := targetDir + "/" + strconv.Itoa(zLevel) + ".png"

	var targImg *image.RGBA

	// check if already exists, if so, keep working on the existing image.
	if Exists(targFileName) {
		imgfile, _ := os.Open(targFileName)
		timg, _, _ := image.Decode(imgfile)
		targImg = CloneToRGBA(timg)
	} else {
		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}
		targImg = image.NewRGBA(image.Rectangle{upLeft, lowRight})
	}

	for y, origImg := range imgs {
		// Set color for each pixel.
		for x := 0; x < width; x++ {
			color := origImg.At(x, zLevel)
			targImg.Set(x, (y + startYIndex), color)
		}
	}
	f, _ := os.Create(targFileName)
	png.Encode(f, targImg)
}

// GenerateSagCutOnX
func GenerateSagCutOnX(xLevel int, startYIndex int, targetDir string, imgs []image.Image, finalHeight int) {
	width := imgs[0].Bounds().Dy()
	height := finalHeight

	targFileName := targetDir + "/" + strconv.Itoa(xLevel) + ".png"

	var targImg *image.RGBA

	// check if already exists, if so, keep working on the existing image.
	if Exists(targFileName) {
		imgfile, _ := os.Open(targFileName)
		timg, _, _ := image.Decode(imgfile)
		targImg = CloneToRGBA(timg)
	} else {
		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}
		targImg = image.NewRGBA(image.Rectangle{upLeft, lowRight})
	}

	for i, origImg := range imgs {
		// Set color for each pixel.
		for y := 0; y < width; y++ {
			color := origImg.At(xLevel, y)
			targImg.Set(y, (i + startYIndex), color)
		}
	}
	f, _ := os.Create(targFileName)
	png.Encode(f, targImg)
}

// DoCoronal - Generates the coronal images from a starting index,
// returns the index where its reached the maximum amount of memory
// allowed to do the processing
func DoCoronal(imgs []image.Image, startIndex int, height int) {
	path := "./coronal-out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if len(imgs) > 0 {

		numCuts := imgs[0].Bounds().Dy() // height
		for i := 0; i < numCuts; i++ {
			GenerateCorCutOnZ(i, startIndex, path, imgs, height)
		}
	}
}

func DoSagital(imgs []image.Image, startIndex int, length int) {
	path := "./sagital-out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if len(imgs) > 0 {

		numCuts := imgs[0].Bounds().Dx()
		for i := 0; i < numCuts; i++ {
			GenerateSagCutOnX(i, startIndex, path, imgs, length)
		}
	}
}

func LoadFiles(startIndex int, files []os.FileInfo, path string, maxMem int) ([]image.Image, int) {
	var imgs []image.Image
	currIndex := 0
	memCost := 0
	for y, f := range files {
		fullPathFile := f.Name()
		fileExt := filepath.Ext(fullPathFile)
		if fileExt != ".png" { // only pngs
			continue
		}

		if memCost/1000000 > maxMem {
			break
		}

		if y > startIndex {
			imgfile, err := os.Open(path + "/" + fullPathFile)
			if err != nil {
				fmt.Println("Image file couldn't be opened.")
				imgfile.Close()
				continue
			}

			img, _, _ := image.Decode(imgfile)
			fi, _ := imgfile.Stat()
			memCost += (int)(fi.Size())
			imgs = append(imgs, img)
			currIndex = y
			imgfile.Close()
		}
	}
	return imgs, currIndex + 1
}

func main() {

	var maxMem int

	argsWithoutProg := os.Args[1:]

	// Get the parameters
	if len(argsWithoutProg) < 1 {
		printUsage()
	}

	if len(argsWithoutProg) > 2 {
		maxMem, _ = strconv.Atoi(argsWithoutProg[2])
	} else {
		maxMem = 100
	}

	counter, theFiles := GetFilesFromExtension("png", argsWithoutProg[0])

	if counter > 0 {
		tempI := 0
		i := 0
		processType := argsWithoutProg[1]
		var imgs []image.Image

		for i < counter {

			imgs, tempI = LoadFiles(tempI, theFiles, argsWithoutProg[0], maxMem)
			switch processType {
			case "coronal":
				DoCoronal(imgs, i, len(theFiles))
				i = tempI
			case "sagital":
				DoSagital(imgs, i, len(theFiles))
				i = tempI
			}
		}
	}
}
