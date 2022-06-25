package model

import _ "encoding/json"

type BitmapTextFontJson3 struct {
	ExportType  string         `json:"ExportType"`
	ExportValue ExportValueObj `json:"ExportValueObj"`
}

func (charmap *BitmapTextFontJson3) GetExportValue() ExportValueObj {
	return charmap.ExportValue
}
