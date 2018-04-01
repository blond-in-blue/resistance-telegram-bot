package main


// Command is a type of bot action that can be executed on a message
type BotCommand struct {
	Matcher func(msg Update) bool
	Execute func(x TeleBot, y Update, w chan BotResponse)
	Name string
	Description string
}

func NewBotCommand(matcher func(Update) bool, execute func(TeleBot,Update, chan BotResponse) , name string, description string) *BotCommand {
	p := new(BotCommand)
	p.Matcher = matcher
	p.Execute = execute
	p.Name = name
	p.Description = description
	return p
}