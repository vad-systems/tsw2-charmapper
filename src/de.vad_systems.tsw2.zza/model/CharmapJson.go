package model

import _ "encoding/json"

type CharEntry struct {
	CharCode string `json:"CharCode"`
	RawDataIndex uint `json:"RawDataIndex"`
	SizeX uint8 `json:"SizeX"`
	SizeY uint8 `json:"SizeY"`
	IsValid bool `json:"bIsValid"`
}

type CharmapJson struct {
	SingleChars []CharEntry `json:"SingleChars"`
	SingleCharIndices []int `json:"SingleCharIndices"`
	MultiChars []CharEntry `json:"MultiChars"`
}
