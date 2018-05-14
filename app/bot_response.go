package main

// BotResponse set by bot when a command is satisfied
type BotResponse struct {
	Text     string
	Pid      string
	Sid      string
	FilePath string
	ChatID   int64
}

func (res BotResponse) GetChatID() int64 {
	return res.ChatID
}

func (res BotResponse) IsTextMessage() bool {
	return res.Text != ""
}

func (res BotResponse) IsPicture() bool {
	return res.Pid != ""
}

func (res BotResponse) IsSticker() bool {
	return res.Sid != ""
}

func (res BotResponse) IsFile() bool {
	return res.FilePath != ""
}

func (res BotResponse) GetTextMessage() string {
	return res.Text
}

func (res BotResponse) GetPicture() string {
	return res.Pid
}

func (res BotResponse) GetSticker() string {
	return res.Sid
}

func (res BotResponse) GetFilePath() string {
	return res.FilePath
}

func NewTextBotResponse(msg string, chatID int64) *BotResponse {
	p := new(BotResponse)
	p.Text = msg
	p.ChatID = chatID
	return p
}

func NewPictureBotResponse(pid string, chatID int64) *BotResponse {
	p := new(BotResponse)
	p.Pid = pid
	p.ChatID = chatID
	return p
}

func NewStickerBotResponse(sid string, chatID int64) *BotResponse {
	p := new(BotResponse)
	p.Sid = sid
	p.ChatID = chatID
	return p
}

func NewFileBotResponse(filePath string, chatID int64) *BotResponse {
	p := new(BotResponse)
	p.FilePath = filePath
	p.ChatID = chatID
	return p
}
