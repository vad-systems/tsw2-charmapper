package main

import (
	"./model"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

var zZAFace font.Face

func initialize() {
	var err error
	zZAFontFileContents, err := os.ReadFile("./zzafont.otf")
	panicOnErr(err)
	zZAFont, err := truetype.Parse(zZAFontFileContents)
	panicOnErr(err)
	zZAFace = truetype.NewFace(zZAFont, &truetype.Options{ Size: 18, SubPixelsX: 2, SubPixelsY: 2})
}

func drawText(img *image.Gray, text string) {
	col := color.Gray{ Y: 0 }
	point := fixed.Point26_6{ X: -32, Y: (21 - 5) * 64 }

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: zZAFace,
		Dot:  point,
	}
	d.DrawString(text)
}

func main() {
	initialize()
	var err error
	charmapFileLoc := flag.String("charmap", "charmap.json", "Location of the charmap.json to use")

	contents, err := ioutil.ReadFile(*charmapFileLoc)
	panicOnErr(err)

	var charmap model.CharmapJson
	err = json.Unmarshal(contents, &charmap)
	panicOnErr(err)
	matrix := ""
	for i := 0 ; i < len(charmap.SingleChars); i++ {
		f, err := os.Create("glyphs/s" + strconv.Itoa(i) + ".png")
		panicOnErr(err)
		img := image.NewGray(image.Rect(0, 0, int(charmap.SingleChars[i].SizeX), int(charmap.SingleChars[i].SizeY)))
		for x := 0 ; x < img.Rect.Dx(); x++ {
			for y := 0 ; y < img.Rect.Dy() ; y++ {
				img.Set(x, y, color.Gray{Y: 255})
			}
		}
		drawText(img, charmap.SingleChars[i].CharCode)
		for x := 0; x < int(charmap.SingleChars[i].SizeX); x++ {
			for y := 0; y < int(charmap.SingleChars[i].SizeY); y++ {
				if r, _, _, a := img.At(x, y).RGBA(); r < a/2 {
					img.Set(x, y, color.Gray{Y: 0})
					matrix += "1"
				} else {
					img.Set(x, y, color.Gray{Y: 255})
					matrix += "0"
				}
			}
		}
		err = png.Encode(f, img)
		panicOnErr(err)
		err = f.Close()
		panicOnErr(err)
	}

	search := ""
	minDiff := 999999999
	offset := 0
	targetX := 0
	for maxX := 7 ; maxX <= 13 ; maxX++ {
		img := image.NewGray(image.Rect(0, 0, maxX, 21))
		for x := 0 ; x < img.Rect.Dx(); x++ {
			for y := 0 ; y < img.Rect.Dy() ; y++ {
				img.Set(x, y, color.Gray{ Y: 255 })
			}
		}
		drawText(img, "y")

		for x := 0 ; x < maxX; x++ {
			for y := 0 ; y < 21 ; y++ {
				if r, _, _, a := img.At(x, y).RGBA(); r < a/2 {
					search += "1"
				} else {
					search += "0"
				}
			}
		}

		for i := 0 ; i < len(matrix) - len(search) ; i++ {
			diff := 0
			for k := 0 ; k < len(search) ; k++ {
				if matrix[i+k] != search[k] {
					diff++
				}
			}
			if diff < minDiff {
				minDiff = diff
				offset = i
				targetX = maxX
			}
		}

		if targetX == maxX {
			fmt.Println("Found potential at Offset =", offset, "with SizeX =", targetX, "and diff of", minDiff)
		}
	}
}
