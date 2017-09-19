package main

// Sticker contains information about a sticker.
type Sticker struct {
	FileID    string     `json:"file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Thumbnail *PhotoSize `json:"thumb"`     // optional
	Emoji     string     `json:"emoji"`     // optional
	FileSize  int        `json:"file_size"` // optional
}
