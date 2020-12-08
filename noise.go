package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
)

//type Pixel
type Pixel struct {
	R int
	G int
	B int
	A int
}

var height, width = 0, 0
var imgLoaded [][]Pixel

// crée une matrice à partir d'une image
// récupère les valeurs RGBA de chaque pixel
func getImg(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height = bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		imgLoaded = append(imgLoaded, row)
	}
	return imgLoaded, nil
}

//conversion uint32 -> uint8
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func encode() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(getRed(pixels[y][x])),
				G: uint8(getGreen(pixels[y][x])),
				B: uint8(getBlue(pixels[y][x])),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func getRed(lePixel Pixel) int {
	return lePixel.R
}

func getGreen(lePixel Pixel) int {
	return lePixel.G
}

func getBlue(lePixel Pixel) int {
	return lePixel.B
}

func getAlpha(lePixel Pixel) int {
	return lePixel.A
}

func main() {

	file, err := os.Open("Tapping-Noise-in-Attic-at-Night.png")

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getImg(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	mean(pixels)
	toBlackAndWhite(pixels)
	encode()

}

func mean(pixels Pixel) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8((getRed(pixels[y+1][x]) + getRed(pixels[y-1][x]) + getRed(pixels[y][x-1]) + getRed(pixels[y][x+1])/4)),
				G: uint8((getGreen(pixels[y+1][x]) + getGreen(pixels[y-1][x]) + getGreen(pixels[y][x-1]) + getGreen(pixels[y][x+1])/4)),
				B: uint8((getBlue(pixels[y+1][x]) + getBlue(pixels[y-1][x]) + getBlue(pixels[y][x-1]) + getBlue(pixels[y][x+1])/4)),
				A: uint8((getAlpha(pixels[y+1][x]) + getAlpha(pixels[y-1][x]) + getAlpha(pixels[y][x-1]) + getAlpha(pixels[y][x+1])/4)),
			})
		}
	}
}

func toBlackAndWhite(img [][]Pixel) {
	for i := 0; i < len(img); i++ {
		for j := 0; j < len(img[i]); j++ {
			pixel := img[i][j]

			newRed := (getRed(pixel) + getGreen(pixel) + getBlue(pixel)) / 3
			newGreen := newRed
			newBlue := newRed

			imgLoaded[i][j] = Pixel{newRed, newGreen, newBlue, 255}
		}
	}
}