package model

type ExportValueObj struct {
	MaxGlyphHeight    int         `json:"MaxGlyphHeight"`
	MaxGlyphWidth     int         `json:"MaxGlyphWidth"`
	GlyphSpacingX     int         `json:"GlyphSpacingX"`
	BCachedData       bool        `json:"bCachedData"`
	CachedDataVersion int         `json:"CachedDataVersion"`
	RawTextureData    []int       `json:"RawTextureData"`
	SingleChars       []CharEntry `json:"SingleChars"`
	SingleCharIndices []int       `json:"SingleCharIndices"`
	MultiChars        []CharEntry `json:"MultiChars"`
}
