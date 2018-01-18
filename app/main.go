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

// Builds and returns commands with url.
func getCommands(telebot Telegram, redditSession RedditAccount, errorLogger func(string), msgBuffer map[string]MessageStack) []func(Update) {

	chatAliases := make(map[string]string)

	return []func(update Update){

		// alias command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "alias-set")
			if matches && commands != "" {

				_, alreadyExists := chatAliases[commands]
				if alreadyExists {
					go telebot.SendMessage(fmt.Sprintf("Someone has already taken the alias '%s'", commands), update.Message.Chat.ID)
				} else {
					chatAliases[commands] = strconv.FormatInt(update.Message.Chat.ID, 10)
					go telebot.SendMessage(fmt.Sprintf("Alias set as: '%s'", commands), update.Message.Chat.ID)
				}

			}
		},

		// alias command
		// func(update Update) {
		// 	matches, _ := getContentFromCommand(update.Message.Text, "alias-get")
		// 	if matches {
		// 		alias, prs := chatAliases[strconv.FormatInt(update.Message.Chat.ID, 10)]
		// 		if prs {
		// 			go telebot.SendMessage("Alias set as: "+alias, update.Message.Chat.ID)
		// 		} else {
		// 			go telebot.SendMessage("No alias set! Use /alias-set <name> to set", update.Message.Chat.ID)
		// 		}
		// 	}
		// },

		// Eli is a furry command
		func(update Update) {
			msg := strings.ToLower(update.Message.Text)
			if (strings.Contains(msg, "eli") || strings.Contains(msg, "b02s2")) && strings.Contains(msg, "furry") {
				go telebot.SendMessage("Actually, "+update.Message.From.UserName+" is the furry", update.Message.Chat.ID)
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
					location := strconv.FormatInt(update.Message.Chat.ID, 10)
					if otherChatID != "" {
						location = otherChatID
						alias, prs := chatAliases[location]
						if prs {
							location = alias
						}
					}
					msgBuffer[location] = msgBuffer[location].Push(*update.Message.ReplyToMessage)
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
				} else {
					go telebot.SendMessage("Reply to message to edge", update.Message.Chat.ID)
				}

			}
		},

		func(update Update) {
			matches, _ := getContentFromCommand(update.Message.Text, "ejaculate")
			if matches {
				msgSent := false
				buffer := msgBuffer[strconv.FormatInt(update.Message.Chat.ID, 10)]
				for msg := range buffer.Everything() {
					msgSent = true
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
				if msgSent == false {
					telebot.SendMessage("Im all tapped out", update.Message.Chat.ID)
				}
				msgBuffer[strconv.FormatInt(update.Message.Chat.ID, 10)] = make([]Message, 0)
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
			if update.Message.Text == "/traps" {
				go telebot.SendMessage("https://www.youtube.com/watch?v=9E1YYSZ9qrk", update.Message.Chat.ID)
			}
		},

		// God command
		func(update Update) {
			if update.Message.Text == "/gg" {
				go telebot.SendMessage("GOD IS GREAT", update.Message.Chat.ID)
			}
		},

		// Rule34 command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "rule34")
			if matches && commands != "" {
				go rule34Search(commands, telebot, update, errorLogger, redditSession)
			}
		},

		// Hedgehog
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "hedgehog")
			if matches && commands != "" {
				go hedgeHogCommand(commands, telebot, update, errorLogger, redditSession)
			}
		},

		// Save command
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "save")
			if matches {
				if commands != "" {
					go SaveCommand(commands, telebot, update, errorLogger, redditSession)
				} else {
					go telebot.SendMessage("Please provide a title for the post.", update.Message.Chat.ID)
				}

			}
		},

		//pokedexSerach
		func(update Update) {
			matches, commands := getContentFromCommand(update.Message.Text, "pokedex")
			if matches && commands != "" {
				go pokedexSearch(commands, telebot, update, errorLogger)
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
					errorLogger("unable to load image: " + err.Error())
					return
				}
				dc := gg.NewContextForImage(im)

				dc.SetRGB(1, 1, 1)
				font, err := truetype.Parse(goregular.TTF)
				if err != nil {
					errorLogger(err.Error())
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
	}
}

// Create our routes
func initRoutes(router *gin.Engine, errors *[]string, msgBuffer map[string]MessageStack) {

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

	timeStarted := getTime()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"restarted": timeStarted,
			"errors":    errors,
		})
	})

	router.GET("/edge/:chatID", func(c *gin.Context) {
		chatID := c.Param("chatID")
		msgs, _ := msgBuffer[chatID]
		//msgs[0].Photo
		log.Println(msgs)
		c.HTML(http.StatusOK, "edge.tmpl", gin.H{
			"messages": msgs,
		})
	})

	router.StaticFS("/media", http.Dir("media"))

}

func listenForUpdates(telebot Telegram, errorLogger func(string), msgBuffer map[string]MessageStack) {

	redditUser := logginToReddit(errorLogger)
	commands := getCommands(telebot, redditUser, errorLogger, msgBuffer)

	for {
		// Sleep first, so if we error out and continue to the next loop, we still end up waiting
		time.Sleep(time.Second)

		updates, err := telebot.GetUpdates()

		if err != nil {
			errorLogger("Error getting updates from telegram: " + err.Error())
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

// Format the current time
func getTime() string {
	t := time.Now()
	return t.Format("Mon Jan _2 15:04:05 UTC-01:00 2006")
}

func logginToReddit(errorLogger func(string)) RedditAccount {

	log.Printf("Logging into: %s\n", os.Getenv("REDDIT_USERNAME"))
	user, err := LoginToReddit(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		"Resistance Telegram Botter",
	)
	if err != nil {
		errorLogger("Error logging into reddit! " + err.Error())
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

	// Create buffers for all the chats.
	allChatBuffers := make(map[string]MessageStack)

	errorMessages := []string{}
	var errorLogger = func(msg string) {
		log.Println(msg)
		newMsg := [...]string{getTime() + ": " + msg}
		errorMessages = append(newMsg[:], errorMessages...)
	}

	teleBot := NewTelegramBot(os.Getenv("TELE_KEY"), errorLogger)

	go listenForUpdates(*teleBot, errorLogger, allChatBuffers)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, &errorMessages, allChatBuffers)
	r.Run(":" + port)

}
