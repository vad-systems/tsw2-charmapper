package main

import (
	"./model"
	"encoding/json"
	"errors"
	"flag"
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

var charmapFileLoc *string

func main() {
	mode := flag.String("mode", "extract", "(extract|stitch|inplace) Mode to run in")
	charmapFileLoc = flag.String("font-json", "BitmapTextFont.json", "Location of the BitmapTextFont.json to use")
	flag.Parse()

	switch *mode {
	case "extract":
		extract()
	case "stitch":
		stitch()
	case "inplace":
		inplace()
	}
}

func extract() {
	var err error

	contents, err := ioutil.ReadFile(*charmapFileLoc)
	panicOnErr(err)

	var charmap []model.Charmap
	var charmap_v3 []model.BitmapTextFontJson3
	err = json.Unmarshal(contents, &charmap_v3)
	if len(charmap_v3) > 0 && len(charmap_v3[0].ExportType) > 0 {
		charmap = make([]model.Charmap, len(charmap_v3))
		for k, d := range charmap_v3 {
			charmap[k] = &d
		}
	} else {
		var charmap_v4 []model.BitmapTextFontJson
		err = json.Unmarshal(contents, &charmap_v4)

		if len(charmap_v4) > 0 && len(charmap_v4[0].Name) > 0 {
			charmap = make([]model.Charmap, len(charmap_v4))
			for k, d := range charmap_v4 {
				charmap[k] = &d
			}
		}
	}
	var e error
	if err != nil {
		e = err
	} else if charmap == nil {
		e = errors.New("No charmap")
	}
	panicOnErr(e)

	for i := range charmap {
		font := charmap[i].GetExportValue()
		dir := "./font" + strconv.Itoa(i)

		if _, err := os.Open(dir); os.IsNotExist(err) {
			err := os.Mkdir(dir, 0755)
			panicOnErr(err)
		}

		for c := range font.SingleChars {
			char := font.SingleChars[c]

			img := image.NewGray(image.Rect(0, 0, int(char.SizeX), int(char.SizeY)))
			for x := 0; x < img.Stride; x++ {
				for y := 0; y < img.Rect.Dy(); y++ {
					if val := uint8(font.RawTextureData[int(char.RawDataIndex)+y*img.Stride+x]); val > 127 {
						img.Set(x, y, color.Gray{Y: 255})
					} else {
						img.Set(x, y, color.Gray{Y: 0})
					}
				}
			}

			f, err := os.Create(dir + "/single_" + strconv.Itoa(c) + ".png")
			panicOnErr(err)
			err = png.Encode(f, img)
			panicOnErr(err)
			err = f.Close()
			panicOnErr(err)
		}

		for c := range font.MultiChars {
			char := font.MultiChars[c]

			img := image.NewGray(image.Rect(0, 0, int(char.SizeX), int(char.SizeY)))
			for x := 0; x < img.Stride; x++ {
				for y := 0; y < img.Rect.Dy(); y++ {
					if val := uint8(font.RawTextureData[int(char.RawDataIndex)+y*img.Stride+x]); val > 127 {
						img.Set(x, y, color.Gray{Y: 255})
					} else {
						img.Set(x, y, color.Gray{Y: 0})
					}
				}
			}

			f, err := os.Create(dir + "/multi_" + strconv.Itoa(c) + ".png")
			panicOnErr(err)
			err = png.Encode(f, img)
			panicOnErr(err)
			err = f.Close()
			panicOnErr(err)
		}
	}
}

func stitch() {
	var err error

	contents, err := ioutil.ReadFile(*charmapFileLoc)
	panicOnErr(err)

	var charmap []model.Charmap
	var charmap_v3 []model.BitmapTextFontJson3
	err = json.Unmarshal(contents, &charmap_v3)
	if len(charmap_v3) > 0 && len(charmap_v3[0].ExportType) > 0 {
		charmap = make([]model.Charmap, len(charmap_v3))
		for k, d := range charmap_v3 {
			charmap[k] = &d
		}
	} else {
		var charmap_v4 []model.BitmapTextFontJson
		err = json.Unmarshal(contents, &charmap_v4)

		if len(charmap_v4) > 0 && len(charmap_v4[0].Name) > 0 {
			charmap = make([]model.Charmap, len(charmap_v4))
			for k, d := range charmap_v4 {
				charmap[k] = &d
			}
		}
	}
	var e error
	if err != nil {
		e = err
	} else if charmap == nil {
		e = errors.New("No charmap")
	}
	panicOnErr(e)

	for i := range charmap {
		font := charmap[i].GetExportValue()
		dir := "./font" + strconv.Itoa(i)

		// We read SingleCharIndices and stitch the files accordingly
		// then provide the new RawTextureData as well as the new SingleChars entries
		charsData := make([]int, 0)
		for c := range font.SingleCharIndices {
			idx := &font.SingleCharIndices[c]
			if *idx == -1 {
				continue
			}

			char := &font.SingleChars[*idx]
			f, err := os.Open(dir + "/single_" + strconv.Itoa(*idx) + ".png")
			panicOnErr(err)
			img, err := png.Decode(f)
			panicOnErr(err)
			char.RawDataIndex = uint(len(charsData))
			char.SizeX = uint8(img.Bounds().Dx())
			char.SizeY = uint8(img.Bounds().Dy())
			for y := 0; y < int(char.SizeY); y++ {
				for x := 0; x < int(char.SizeX); x++ {
					r, _, _, a := img.At(x, y).RGBA()
					charsData = append(charsData, int((r/a)*255))
				}
			}
		}

		for c := range font.MultiChars {
			char := &font.MultiChars[c]
			f, err := os.Open(dir + "/multi_" + strconv.Itoa(c) + ".png")
			panicOnErr(err)
			img, err := png.Decode(f)
			panicOnErr(err)
			char.RawDataIndex = uint(len(charsData))
			char.SizeX = uint8(img.Bounds().Dx())
			char.SizeY = uint8(img.Bounds().Dy())
			for y := 0; y < int(char.SizeY); y++ {
				for x := 0; x < int(char.SizeX); x++ {
					r, _, _, a := img.At(x, y).RGBA()
					charsData = append(charsData, int((r/a)*255))
				}
			}
		}

		font.RawTextureData = charsData

		rawData := make([]uint8, 0)
		for b := range font.RawTextureData {
			rawData = append(rawData, uint8(font.RawTextureData[b]))
		}

		err := os.WriteFile("BitmapTextFont."+strconv.Itoa(i)+".texture", rawData, 0644)
		newJson, err := json.Marshal(font.RawTextureData)
		err = ioutil.WriteFile("BitmapTextFont."+strconv.Itoa(i)+".texture.json", newJson, 0644)
		panicOnErr(err)
	}

	newJson, err := json.Marshal(charmap)
	panicOnErr(err)

	err = ioutil.WriteFile("BitmapTextFont.out.json", newJson, 0644)
	panicOnErr(err)
}

func inplace() {
	var err error

	contents, err := ioutil.ReadFile(*charmapFileLoc)
	panicOnErr(err)

	var charmap []model.Charmap
	var charmap_v3 []model.BitmapTextFontJson3
	err = json.Unmarshal(contents, &charmap_v3)
	if len(charmap_v3) > 0 && len(charmap_v3[0].ExportType) > 0 {
		charmap = make([]model.Charmap, len(charmap_v3))
		for k, d := range charmap_v3 {
			charmap[k] = &d
		}
	} else {
		var charmap_v4 []model.BitmapTextFontJson
		err = json.Unmarshal(contents, &charmap_v4)

		if len(charmap_v4) > 0 && len(charmap_v4[0].Name) > 0 {
			charmap = make([]model.Charmap, len(charmap_v4))
			for k, d := range charmap_v4 {
				charmap[k] = &d
			}
		}
	}
	var e error
	if err != nil {
		e = err
	} else if charmap == nil {
		e = errors.New("No charmap")
	}
	panicOnErr(e)

	for i := range charmap {
		font := charmap[i].GetExportValue()
		dir := "./font" + strconv.Itoa(i)

		for c := range font.SingleCharIndices {
			idx := &font.SingleCharIndices[c]
			if *idx == -1 {
				continue
			}

			char := &font.SingleChars[*idx]
			f, err := os.Open(dir + "/single_" + strconv.Itoa(*idx) + ".png")
			panicOnErr(err)
			img, err := png.Decode(f)
			panicOnErr(err)
			charData := make([]int, 0)
			for y := 0; y < int(char.SizeY); y++ {
				for x := 0; x < int(char.SizeX); x++ {
					r, _, _, a := img.At(x, y).RGBA()
					charData = append(charData, int((r/a)*255))
				}
			}
			for len(font.RawTextureData) < int(char.RawDataIndex)+len(charData) {
				font.RawTextureData = append(font.RawTextureData, 0)
			}

			for b := range charData {
				font.RawTextureData[int(char.RawDataIndex)+b] = charData[b]
			}
		}

		for c := range font.MultiChars {
			char := &font.MultiChars[c]
			f, err := os.Open(dir + "/multi_" + strconv.Itoa(c) + ".png")
			panicOnErr(err)
			img, err := png.Decode(f)
			panicOnErr(err)
			charData := make([]int, 0)
			for y := 0; y < int(char.SizeY); y++ {
				for x := 0; x < int(char.SizeX); x++ {
					r, _, _, a := img.At(x, y).RGBA()
					charData = append(charData, int((r/a)*255))
				}
			}
			for len(font.RawTextureData) < int(char.RawDataIndex)+len(charData) {
				font.RawTextureData = append(font.RawTextureData, 0)
			}

			for b := range charData {
				font.RawTextureData[int(char.RawDataIndex)+b] = charData[b]
			}
		}

		rawData := make([]uint8, 0)
		for b := range font.RawTextureData {
			rawData = append(rawData, uint8(font.RawTextureData[b]))
		}

		err := os.WriteFile("BitmapTextFont."+strconv.Itoa(i)+".texture", rawData, 0644)
		panicOnErr(err)
		newJson, err := json.Marshal(font.RawTextureData)
		err = ioutil.WriteFile("BitmapTextFont."+strconv.Itoa(i)+".texture.json", newJson, 0644)
	}

	newJson, err := json.Marshal(charmap)
	panicOnErr(err)

	err = ioutil.WriteFile("BitmapTextFont.out.json", newJson, 0644)
	panicOnErr(err)
}
