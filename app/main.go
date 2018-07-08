// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	// "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func messageContainsCommandMatcher(command string) func(Update) bool {
	return func(update Update) bool {
		return messageContainsCommand(update.Message.Text, command)
	}
}

func messageContainsCommand(message string, command string) bool {
	return len(strings.SplitAfter(message, fmt.Sprintf("/%s", command))) > 1
}

func getContentFromCommand(message string, command string) string {
	commands := strings.SplitAfter(message, fmt.Sprintf("/%s", command))
	if len(commands) > 1 {
		return strings.TrimSpace(commands[1])
	}
	return ""
}

// Builds and returns commands
func getCommands() []BotCommand {

	return []BotCommand{

		BotCommand{
			Name:        "Help",
			Description: "list of commands",
			Matcher: func(update Update) bool {
				return update.Message.Text == "/help"
			},
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				var returnMsg bytes.Buffer
				returnMsg.WriteString("COMMANDS\n ")
				for _, command := range bot.GetCommands() {
					returnMsg.WriteString(fmt.Sprintf("\n<b>%s</b> - %s\n", command.Name, command.Description))
				}
				respChan <- *NewTextBotResponse(returnMsg.String(), update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "traps",
			Description: "just a friendly reminder",
			Matcher:     messageContainsCommandMatcher("traps"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("https://www.youtube.com/watch?v=9E1YYSZ9qrk", update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "ping",
			Description: "check if the bot is listening",
			Matcher:     messageContainsCommandMatcher("ping"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("fuck you want?", update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "wuh",
			Description: "stop fucking stop",
			Matcher:     messageContainsCommandMatcher("wuh"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("https://www.youtube.com/watch?v=j3z7vjs1E18", update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "Ahem",
			Description: "You are a furry arn't you",
			Matcher: func(update Update) bool {
				return update.Message.Text == "ahem" && update.Message.ReplyToMessage != nil
			},
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s is actually the furry", update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "Alias Set",
			Description: "Alias the chat for other commands like edge, <code>/alias-set resistance</code>",
			Matcher:     messageContainsCommandMatcher("alias-set"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				alias := getContentFromCommand(update.Message.Text, "alias-set")

				if alias != "" {
					_, alreadyExists := bot.IsAliasSet(alias)
					if alreadyExists {
						respChan <- *NewTextBotResponse(fmt.Sprintf("Someone has already taken the alias '%s'", alias), update.Message.Chat.ID)
					} else {
						bot.SetChatAlias(alias, update.Message.Chat.ID)
						respChan <- *NewTextBotResponse(fmt.Sprintf("Alias set as: '%s'", alias), update.Message.Chat.ID)
					}
				}
			},
		},

		BotCommand{
			Name:        "Password",
			Description: "Gives you chat id for edged site, /password",
			Matcher:     messageContainsCommandMatcher("password"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse(strconv.FormatInt(update.Message.Chat.ID, 10), update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "Leaving",
			Description: "Cause you want more attention, /leaving",
			Matcher:     messageContainsCommandMatcher("leaving"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				respChan <- *NewTextBotResponse("Y’all are miserable people who demand to be right at all times, even when you have no experience on the subject. I’m sick of it and it’s gonna be better to just not have to deal with it. So I’m out. You’re mostly all hugely negative impacts on mine and others’ lives. Obviously some of you I still consider friends, but it’s really gotten to where the animosity of a few people make this little group irredeemably shitty for me to be a part of. Or irredeemably shitty in general really. There’s really just no good part of Resistance. ", update.Message.Chat.ID)
				respChan <- *NewTextBotResponse("Especially when it’s pretty much exclusively me that seems to be the target of all the hate. It just feels super mean spirited. It would be different if you attacked everybody the same way, but you don’t. And don’t even try to fucking pretend you do. I’ve genuinely never felt liked here. ", update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "Edge",
			Description: "Hide messages for later, reply to a message with /edge",
			Matcher: func(update Update) bool {
				return update.Message.ReplyToMessage != nil && update.Message.Text == "/edge"
			},
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				bot.PushMessageToChatBuffer(strconv.FormatInt(update.Message.Chat.ID, 10), *update.Message.ReplyToMessage)
				if update.Message.ReplyToMessage.Photo != nil {
					photos := *update.Message.ReplyToMessage.Photo
					go bot.GetFile(photos[0].FileID)
				}

				if update.Message.ReplyToMessage.Sticker != nil {
					sticker := *update.Message.ReplyToMessage.Sticker
					go bot.GetFile(sticker.FileID)
				}

				go bot.deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
				go bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s edged %s", update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
			},
		},

		BotCommand{
			Name:        "Ejaculate",
			Description: "Release all the messages that have been edged with /ejaculate",
			Matcher: func(update Update) bool {
				return update.Message.Text == "/ejaculate"
			},
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				msgSentCount := 0
				buffer := bot.ClearBuffer(update.Message.Chat.ID)
				for msg := range buffer.Everything() {
					msgSentCount++
					if msg.Photo != nil {
						photos := *msg.Photo
						respChan <- *NewTextBotResponse(msg.From.UserName+" sent:", update.Message.Chat.ID)
						respChan <- *NewPictureBotResponse(photos[0].FileID, update.Message.Chat.ID)
					} else if msg.Sticker != nil {
						respChan <- *NewTextBotResponse(msg.From.UserName+" sent:", update.Message.Chat.ID)
						respChan <- *NewStickerBotResponse(msg.Sticker.FileID, update.Message.Chat.ID)
					} else if msg.Document != nil {
						respChan <- *NewTextBotResponse(msg.From.UserName+" sent:", update.Message.Chat.ID)
						respChan <- *NewStickerBotResponse(msg.Document.FileID, update.Message.Chat.ID)
					} else {
						respChan <- *NewTextBotResponse(fmt.Sprintf("%s sent:\n%s", msg.From.UserName, msg.Text), update.Message.Chat.ID)
					}
				}

				if msgSentCount == 0 {
					respChan <- *NewTextBotResponse("I'm not usually like this. Maybe if you do something sexy it'll start working", update.Message.Chat.ID)
				} else if msgSentCount < 5 {
					respChan <- *NewTextBotResponse("Normally I'm not that quick", update.Message.Chat.ID)
				} else if msgSentCount < 10 {
					respChan <- *NewTextBotResponse("I need a ciggarette after that", update.Message.Chat.ID)
				} else {
					respChan <- *NewTextBotResponse("HOLY FUCK I NEEDED THAT, sorry about the mess", update.Message.Chat.ID)
				}
			},
		},

		BotCommand{
			Name:        "Rush",
			Description: "Obliterate your opponent",
			Matcher:     messageContainsCommandMatcher("rush"),
			Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
				wholeCommand := getContentFromCommand(update.Message.Text, "rush")

				if wholeCommand == "" {
					return
				}
				commands := strings.Split(wholeCommand, " with ")

				name := commands[0]
				attack := ""

				if len(commands) > 1 {
					attack = commands[1]
				}
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					bot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 70,
				})

				for i := 0; i < 9; i++ {
					im, err := gg.LoadPNG(fmt.Sprintf("trunks2/F_00%d.png", i))
					if err != nil {
						bot.errorReport.Log("unable to load image: " + err.Error())
						return
					}
					dc := gg.NewContextForImage(im)
					dc.SetRGB(1, 0, 0)
					dc.SetFontFace(face)
					dc.DrawStringAnchored(update.Message.From.UserName, 950, 120, 0.0, 0.0)
					dc.DrawStringAnchored(name, 750, 600, 0.0, 0.0)
					dc.SavePNG(fmt.Sprintf("trunks2out/F_00%d.png", i))
				}

				if attack != "" {
					for i := 34; i < 42; i++ {
						im, err := gg.LoadPNG(fmt.Sprintf("trunks2/F_0%d.png", i))
						if err != nil {
							bot.errorReport.Log("unable to load image: " + err.Error())
							return
						}
						dc := gg.NewContextForImage(im)
						dc.SetRGB(1, 0, 0)
						dc.SetFontFace(face)
						dc.DrawStringAnchored(attack, 250, 300, 0.0, 0.0)
						dc.SavePNG(fmt.Sprintf("trunks2out/F_0%d.png", i))
					}
				}

				StichPicturesTogether("trunks2out")
				respChan <- *NewFileBotResponse("movie.gif", update.Message.Chat.ID)
			},
		},

		holyCommand,
		killCommand,
		rule34Command,
		hedgehogCommand,
		saveCommand,
		pokedexCommand,
		swallowCommand,
		doitCommand,
		repostCommand,
		mockCommand,
	}

}

// Create our routes
func initRoutes(router *gin.Engine, telebot TeleBot) {
	router.SetFuncMap(template.FuncMap{
		"pictureDeref": func(i *[]PhotoSize) PhotoSize {
			if i == nil {
				return PhotoSize{}
			}

			photos := *i
			return photos[0]
		},
		"stickerDeref": func(i *Sticker) Sticker {
			if i == nil {
				return Sticker{}
			}
			return *i
		},
	})

	router.LoadHTMLGlob("templates/*.tmpl")

	timeStarted := GetTime()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"restarted": timeStarted,
			"errors":    telebot.errorReport.Generate(),
		})
	})

	router.GET("/edge/:chatID", func(c *gin.Context) {
		chatID := c.Param("chatID")
		msgs := telebot.ChatBuffer(chatID)
		c.HTML(http.StatusOK, "edge.tmpl", gin.H{
			"messages": msgs,
		})
	})

	router.StaticFS("/media", http.Dir("media"))
}

func listenForUpdates(telebot TeleBot) {

	for {
		// Sleep first, so if we error out and continue to the next loop, we still end up waiting
		time.Sleep(time.Second)

		updates, err := telebot.GetUpdates()

		if err != nil {
			telebot.errorReport.Log("Error getting updates from telegram: " + err.Error())
			continue
		}

		// Dispatch incoming messages to appropriate functions
		for _, update := range updates {
			if update.Message != nil {
				log.Println(update.Message.ToString())
				telebot.OnMessage(update)
			}
		}

	}
}

func logginToReddit(errorReport Report) RedditAccount {

	log.Printf("Logging into: %s\n", os.Getenv("REDDIT_USERNAME"))
	user, err := LoginToReddit(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		"Resistance Telegram Botter",
	)
	if err != nil {
		errorReport.Log("Error logging into reddit! " + err.Error())
	} else {
		log.Println(fmt.Sprintf("Succesfully logged in."))
	}

	return user
}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}
	log.Printf("Starting bot using port %s\n", port)

	errorReport := NewReport()
	redditUser := logginToReddit(*errorReport)
	teleBot := NewTelegramBot(os.Getenv("TELE_KEY"), *errorReport, redditUser, getCommands())
	teleBot.Start()

	go listenForUpdates(*teleBot)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, *teleBot)
	r.Run(":" + port)

}
