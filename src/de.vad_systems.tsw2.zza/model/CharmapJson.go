package model

import _ "encoding/json"

type BitmapTextFontJson struct {
	Type       string         `json:"Type"`
	Name       string         `json:"Name"`
	Properties ExportValueObj `json:"Properties"`
}

func (charmap *BitmapTextFontJson) GetExportValue() ExportValueObj {
	return charmap.Properties
}
