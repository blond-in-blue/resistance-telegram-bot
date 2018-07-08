package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func random(min, max int) int {
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

		if update.Message.ReplyToMessage != nil {
			runes := []rune{}
			newString := ""

			for pos, char := range update.Message.ReplyToMessage.Text {
				if pos != -1 {

				}
				forkInTheRoad := random(0, 2)

				switch forkInTheRoad {
				case 0:
					{
						char = unicode.ToUpper(char)
					}
				case 1:
					{
						char = unicode.ToLower(char)
					}
				}

				runes = append(runes, char)
			}

			for index, rune := range runes {
				if index != -1 {
					newString += joinRunes(rune)
				}
			}

			respChan <- *NewTextBotResponse(fmt.Sprintf("%s", newString), update.Message.Chat.ID)

		} else {
			respChan <- *NewTextBotResponse("reply to a message, retart", update.Message.Chat.ID)
		}
	},
}
