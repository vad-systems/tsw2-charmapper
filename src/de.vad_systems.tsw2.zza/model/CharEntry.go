package model

type CharEntry struct {
	CharCode     string `json:"CharCode"`
	RawDataIndex uint   `json:"RawDataIndex"`
	SizeX        uint8  `json:"SizeX"`
	SizeY        uint8  `json:"SizeY"`
	IsValid      bool   `json:"bIsValid"`
}
