package main

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID          int      `json:"update_id" binding:"required"`
	Message           *Message `json:"message" binding:"required"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
}

// All updates recieved from getUpdate
type BatchUpdates struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
