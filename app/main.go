// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func getContentFromCommand(message string, command string) (bool, string) {
	commands := strings.SplitAfter(message, fmt.Sprintf("/%s", command))
	if len(commands) > 1 {
		return true, strings.TrimSpace(commands[1])
	}
	return false, ""
}

var molyReplacer = strings.NewReplacer("h", "m", "H", "M")

// Builds and returns commands with url.
func getCommands(telebot TeleBot) []func(Update) {

	return []func(update Update){

		// alias command
		func(update Update) {
			matches, alias := getContentFromCommand(update.Message.Text, "alias-set")
			if matches && alias != "" {

				_, alreadyExists := telebot.IsAliasSet(alias)
				if alreadyExists {
					go telebot.SendMessage(fmt.Sprintf("Someone has already taken the alias '%s'", alias), update.Message.Chat.ID)
				} else {
					telebot.SetChatAlias(alias, update.Message.Chat.ID)
					go telebot.SendMessage(fmt.Sprintf("Alias set as: '%s'", alias), update.Message.Chat.ID)
				}

			}
		},

		// Reactions
		func(update Update) {
			for key, value := range reactions {
				matches, _ := getContentFromCommand(update.Message.Text, key)
				if matches {
					go telebot.SendMessage(value, update.Message.Chat.ID)
				}
			}
		},

		func(update Update) {
			if update.Message.Text == "ahem" && update.Message.ReplyToMessage != nil {
				go telebot.SendMessage("Actually, "+update.Message.ReplyToMessage.From.UserName+" is the furry", update.Message.Chat.ID)
			}
		},

		func(update Update) {
			matches, toMatch := getContentFromCommand(update.Message.Text, "s/")
			if matches && toMatch != "" && update.Message.ReplyToMessage != nil {

				aggMessage := update.Message.ReplyToMessage.Text

				for _, line := range strings.SplitAfter(toMatch, "\n") {
					cmds := strings.Split(line, "/")

					if len(cmds) != 2 {
						return
					}

					re, err := regexp.Compile(cmds[0])

					if err == nil {
						aggMessage = re.ReplaceAllString(aggMessage, strings.TrimSpace(cmds[1]))
					} else {
						go telebot.SendMessage(fmt.Sprintf("<b>Invalid expression:</b>\n%s", err.Error()), update.Message.Chat.ID)
					}

				}

				if aggMessage != update.Message.ReplyToMessage.Text {
					go telebot.SendMessage(fmt.Sprintf("<b>Did you mean</b>:\n%s", aggMessage), update.Message.Chat.ID)
				}

			}
		},

		func(update Update) {
			re := regexp.MustCompile("[Hh]+[Oo]+[Ll]+[Yy]+")
			if re.FindString(update.Message.Text) == update.Message.Text && update.Message.Text != "" {
				go telebot.SendMessage(molyReplacer.Replace(update.Message.Text), update.Message.Chat.ID)
			}
		},

		func(update Update) {
			if update.Message.Text == "/password" {
				go telebot.SendMessage(strconv.FormatInt(update.Message.Chat.ID, 10), update.Message.Chat.ID)
			}
		},

		func(update Update) {
			matches, otherChatID := getContentFromCommand(update.Message.Text, "edge")
			if matches {
				if update.Message.ReplyToMessage != nil {
					telebot.PushMessageToChatBuffer(otherChatID, *update.Message.ReplyToMessage)
					if update.Message.ReplyToMessage.Photo != nil {
						photos := *update.Message.ReplyToMessage.Photo
						go telebot.GetImage(photos[0].FileID)
					}
					if update.Message.ReplyToMessage.Sticker != nil {
						sticker := *update.Message.ReplyToMessage.Sticker
						go telebot.GetImage(sticker.FileID)
					}
					go telebot.deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
					go telebot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					go telebot.SendMessage(fmt.Sprintf("%s edged %s", update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
				} else {
					go telebot.SendMessage("Reply to message to edge", update.Message.Chat.ID)
				}

			}
		},

		func(update Update) {
			matches, _ := getContentFromCommand(update.Message.Text, "ejaculate")
			if matches {
				go func() {
					msgSentCount := 0
					for msg := range telebot.ClearBuffer(update.Message.Chat.ID) {
						msgSentCount += 1
						if msg.Photo != nil {
							photos := *msg.Photo
							telebot.SendMessage(msg.From.UserName+" sent:", update.Message.Chat.ID)
							telebot.SendPhotoByID(photos[0].FileID, update.Message.Chat.ID)
						} else if msg.Sticker != nil {
							telebot.SendMessage(msg.From.UserName+" sent:", update.Message.Chat.ID)
							telebot.SendSticker(msg.Sticker.FileID, update.Message.Chat.ID)
						} else {
							telebot.SendMessage(fmt.Sprintf("%s sent:\n%s", msg.From.UserName, msg.Text), update.Message.Chat.ID)
						}
					}

					if msgSentCount == 0 {
						telebot.SendMessage("I'm not usually like this. Maybe if you do something sexy it'll start working", update.Message.Chat.ID)
					} else if msgSentCount < 5 {
						telebot.SendMessage("Normally I'm not that quick", update.Message.Chat.ID)
					} else if msgSentCount < 10 {
						telebot.SendMessage("I need a ciggarette after that", update.Message.Chat.ID)
					} else {
						telebot.SendMessage("HOLY FUCK I NEEDED THAT, sorry about the mess", update.Message.Chat.ID)
					}
				}()
			}
		},

		// Kill command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "kill")
			if matches && commands != "" {
				n := rand.Int() % len(killStatements)
				go telebot.SendMessage(commands+killStatements[n], update.Message.Chat.ID)
			}
		},

		// Traps command
		func(update Update) {
			matches, _ := getContentFromCommand(update.Message.Text, "traps")
			if matches {
				go telebot.SendMessage("https://www.youtube.com/watch?v=9E1YYSZ9qrk", update.Message.Chat.ID)
			}
		},

		// God command
		func(update Update) {
			matches, _ := getContentFromCommand(update.Message.Text, "gg")
			if matches {
				go telebot.SendMessage("GOD IS GREAT", update.Message.Chat.ID)
			}
		},

		// Rule34 command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "rule34")
			if matches && commands != "" {
				go rule34Search(commands, telebot, update)
			}
		},

		// Hedgehog
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "hedgehog")
			if matches && commands != "" {
				go hedgeHogCommand(commands, telebot, update)
			}
		},

		// Ping
		func(update Update) {
			if update.Message.Text == "/ping" {
				go telebot.SendMessage("fuck u want", update.Message.Chat.ID)
			}
		},

		// Save command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "save")
			if matches {
				if commands != "" {
					go SaveCommand(commands, telebot, update)
				} else {
					go telebot.SendMessage("Please provide a title for the post.", update.Message.Chat.ID)
				}

			}
		},

		//pokedexSerach
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "pokedex")
			if matches && commands != "" {
				go pokedexSearch(commands, telebot, update)
			}
		},

		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "murder")
			if matches && commands != "" {

				if commands == "me" {
					commands = update.Message.From.UserName
				}

				im, err := gg.LoadPNG("murder/test.png")
				if err != nil {
					telebot.errorReport.Log("unable to load image: " + err.Error())
					return
				}
				dc := gg.NewContextForImage(im)

				dc.SetRGB(1, 1, 1)
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					telebot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 70,
				})

				dc.SetFontFace(face)
				dc.DrawStringAnchored(commands, 500, 120, 0.0, 0.0)
				dc.SavePNG("media/out.png")
				telebot.sendImage("media/out.png", update.Message.Chat.ID)
			}
		},

		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "valentines")
			if matches && commands != "" {

				if commands == "me" {
					telebot.SendMessage("Stop trying to give yourself love", update.Message.Chat.ID)
					return
				}

				im, err := gg.LoadPNG("murder/valentines.png")
				if err != nil {
					telebot.errorReport.Log("unable to load image: " + err.Error())
					return
				}
				dc := gg.NewContextForImage(im)

				dc.SetRGB(1, 1, 1)
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					telebot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 19,
				})

				dc.SetFontFace(face)
				dc.DrawStringAnchored(commands, 277, 175, 0.0, 0.0)
				dc.DrawStringAnchored(update.Message.From.UserName, 297, 195, 0.0, 0.0)
				dc.SavePNG("media/ValentineOut.png")
				telebot.sendImage("media/ValentineOut.png", update.Message.Chat.ID)
			}
		},

		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "fight")
			if matches && commands != "" {

				fighters := strings.Split(commands, " and ")

				if len(fighters) < 2 {
					return
				}

				left := strings.TrimSpace(fighters[0])
				right := strings.TrimSpace(fighters[1])

				if left == "me" {
					left = update.Message.From.UserName
				}

				if right == "me" {
					right = update.Message.From.UserName
				}

				im, err := gg.LoadPNG("murder/rooster fighting.png")
				if err != nil {
					telebot.errorReport.Log("unable to load image: " + err.Error())
					return
				}
				dc := gg.NewContextForImage(im)

				dc.SetRGB(1, 1, 1)
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					telebot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 70,
				})

				dc.SetFontFace(face)
				dc.DrawStringAnchored(left, 300, 200, 0.0, 0.0)
				dc.DrawStringAnchored(right, 1200, 180, 0.0, 0.0)
				dc.SavePNG("media/roosterOut.png")
				telebot.sendImage("media/roosterOut.png", update.Message.Chat.ID)
			}

		},

		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "hunt")
			if matches && commands != "" {

				fighters := strings.Split(commands, " and ")

				if len(fighters) < 2 {
					return
				}

				left := strings.TrimSpace(fighters[0])
				right := strings.TrimSpace(fighters[1])

				if left == "me" {
					left = update.Message.From.UserName
				}

				if right == "me" {
					right = update.Message.From.UserName
				}

				im, err := gg.LoadPNG("murder/hunt.png")
				if err != nil {
					telebot.errorReport.Log("unable to load image: " + err.Error())
					return
				}
				dc := gg.NewContextForImage(im)

				dc.SetRGB(1, 1, 1)
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					telebot.errorReport.Log(err.Error())
				}
				face := truetype.NewFace(font, &truetype.Options{
					Size: 40,
				})

				dc.SetFontFace(face)
				dc.DrawStringAnchored(left, 100, 220, 0.0, 0.0)
				dc.DrawStringAnchored(right, 300, 375, 0.0, 0.0)
				dc.SavePNG("media/huntOut.png")
				telebot.sendImage("media/huntOut.png", update.Message.Chat.ID)

			}

		},
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

	commands := getCommands(telebot)

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
				for _, command := range commands {
					command(update)
				}
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
	teleBot := NewTelegramBot(os.Getenv("TELE_KEY"), *errorReport, redditUser)

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
