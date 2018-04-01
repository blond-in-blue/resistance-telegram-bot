
package main

import (
	"strings"
	"regexp"
)

var molyReplacer = strings.NewReplacer("h", "m", "H", "M")
var holyRegex = regexp.MustCompile("[Hh]+[Oo]+[Ll]+[Yy]+")

var holyCommand = BotCommand{
	Name: "holy",
	Description: "moly",
	Matcher: func(update Update) bool {
		return holyRegex.FindString(update.Message.Text) == update.Message.Text && update.Message.Text != "" 
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse){
		respChan <- *NewTextBotResponse(molyReplacer.Replace(update.Message.Text), update.Message.Chat.ID)
	},
}