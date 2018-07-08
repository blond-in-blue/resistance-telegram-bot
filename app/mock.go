package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func generateRandomNumberInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func joinRunes(runes ...rune) string {
	var sb strings.Builder
	for _, r := range runes {
		sb.WriteRune(r)
	}
	return sb.String()
}

var mockCommand = BotCommand{
	Name:        "Mock",
	Description: "Reply to a user's message to repeat it in a mocking manner",
	Matcher:     messageContainsCommandMatcher("mock"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {

		inputMessage := ""
		generatedMessage := ""
		generatedMessageRunes := []rune{}

		if update.Message.ReplyToMessage != nil {
			inputMessage = update.Message.ReplyToMessage.Text
		} else if strings.ToLower(update.Message.Text) != "/mock" {
			inputMessage = update.Message.Text[6:len(update.Message.Text)]
		} else {
			inputMessage = "give me something to mock, retart"
		}

		for index, currentRune := range inputMessage {
			if index == -1 {
				return
			}

			caseMode := generateRandomNumberInRange(0, 2)
			switch caseMode {
			case 0:
				{
					currentRune = unicode.ToUpper(currentRune)
				}
			case 1:
				{
					currentRune = unicode.ToLower(currentRune)
				}
			}

			generatedMessageRunes = append(generatedMessageRunes, currentRune)
		}

		for index, rune := range generatedMessageRunes {
			if index == -1 {
				return
			}

			generatedMessage += joinRunes(rune)
		}

		respChan <- *NewTextBotResponse(fmt.Sprintf("%s", generatedMessage), update.Message.Chat.ID)

	},
}
