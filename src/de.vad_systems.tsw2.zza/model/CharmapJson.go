package model

import _ "encoding/json"

type CharEntry struct {
	CharCode string `json:"CharCode"`
	RawDataIndex uint `json:"RawDataIndex"`
	SizeX uint8 `json:"SizeX"`
	SizeY uint8 `json:"SizeY"`
	IsValid bool `json:"bIsValid"`
}

type ExportValueObj struct {
	MaxGlyphHeight int `json:"MaxGlyphHeight"`
	MaxGlyphWidth int `json:"MaxGlyphWidth"`
	GlyphSpacingX int `json:"GlyphSpacingX"`
	BCachedData bool `json:"bCachedData"`
	CachedDataVersion int `json:"CachedDataVersion"`
	RawTextureData []int `json:"RawTextureData"`
	SingleChars []CharEntry `json:"SingleChars"`
	SingleCharIndices []int `json:"SingleCharIndices"`
	MultiChars []CharEntry `json:"MultiChars"`
}

type BitmapTextFontJson struct {
	ExportType string `json:"ExportType"`
	ExportValue ExportValueObj `json:"ExportValue"`
}
