package main

var swallowCommand = BotCommand{
	Name:        "Swallow",
	Description: "Delete all messages in the buffer",
	Matcher:     messageContainsCommandMatcher("swallow"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		buffer := bot.ClearBuffer(update.Message.Chat.ID)
		returnMessage := "No Messages :/"
		if buffer.Size() > 0 {
			returnMessage = "Thanks Daddy :) Messages all gone"
		}
		respChan <- *NewTextBotResponse(returnMessage, update.Message.Chat.ID)
	},
}
