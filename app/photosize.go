package main

// PhotoSize contains information about photos.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"` // optional
}

// GetImageResponse
type GetImageResponse struct {
	Ok     bool                 `json:"ok"`
	Result GetImageResponseData `json:"result"`
}

// GetImageResponseData
type GetImageResponseData struct {
	FilePath string `json:"file_path"`
}
