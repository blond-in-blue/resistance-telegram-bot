package main

// BotResponse set by bot when a command is satisfied
type BotResponse struct {
	Text string
	Pid string
	Sid string
	ChatID int64
}

func (res BotResponse) GetChatID() int64 {
	return res.ChatID
}

func (res BotResponse) IsTextMessage() bool {
	return res.Text != ""
}

func (res BotResponse) GetTextMessage() string {
	return res.Text
}

func  NewTextBotResponse(msg string, chatID int64) *BotResponse{
	p := new(BotResponse)
	p.Text = msg
	p.ChatID = chatID
    return p
}

func  NewPictureBotResponse(pid string) *BotResponse{
	p := new(BotResponse)
    p.Pid = pid
    return p
}

func  NewStickerBotResponse(sid string) *BotResponse{
	p := new(BotResponse)
    p.Sid = sid
    return p
}