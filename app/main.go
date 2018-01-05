// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Builds and returns commands with url.
func getCommands(telebot Telegram, errorLogger func(string), redditSession *http.Cookie, modhash string) []func(Update) {

	return []func(update Update){

		// Eli is a furry command
		func(update Update) {
			if strings.Contains(strings.ToLower(update.Message.Text), "eli") && strings.Contains(strings.ToLower(update.Message.Text), "furry") {
				go telebot.SendMessage("Actually, "+update.Message.From.UserName+" is the furry", update.Message.Chat.ID)
			}

		},

		// Kill command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "/kill")
			if len(commands) > 1 {
				n := rand.Int() % len(killStatements)
				go telebot.SendMessage(strings.TrimSpace(commands[1])+killStatements[n], update.Message.Chat.ID)
			}
		},

		// Traps command
		func(update Update) {
			if strings.Contains(update.Message.Text, "/traps") {
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
			commands := strings.SplitAfter(update.Message.Text, "/rule34")
			if len(commands) > 1 {
				go rule34Search(strings.TrimSpace(commands[1]), telebot, update, errorLogger, redditSession)
			}
		},

		// Save command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "/save")
			if len(commands) > 1 {
				go SaveCommand(strings.TrimSpace(commands[1]), telebot, update, errorLogger, redditSession, modhash)
			}
		},

		//pokedexSerach
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "/pokedex")
			if len(commands) > 1 {
				go pokedexSerach(strings.TrimSpace(commands[1]), telebot, update, errorLogger)
			}
		},

		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "/murder")
			if len(commands) > 1 {
				text := strings.TrimSpace(commands[1])
				// Don't put anything if they didn't give us anything
				if text == "" {
					return
				}

				if text == "me" {
					text = update.Message.From.UserName
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
				dc.DrawStringAnchored(text, 500, 120, 0.0, 0.0)
				dc.SavePNG("media/out.png")
				telebot.sendImage(update.Message.Chat.ID)
			}
		},
	}
}

// Create our routes
func initRoutes(router *gin.Engine, errors *[]string) {

	router.LoadHTMLFiles("templates/index.tmpl")

	timeStarted := getTime()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"restarted": timeStarted,
			"errors":    errors,
		})
	})

	router.StaticFS("/media", http.Dir("media"))

}

func listenForUpdates(telebot Telegram, errorLogger func(string)) {

	cookie, modhash := logginToReddit(errorLogger)
	commands := getCommands(telebot, errorLogger, cookie, modhash)

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
				log.Println("Msg: " + update.Message.Text)
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

func logginToReddit(errorLogger func(string)) (*http.Cookie, string) {

	log.Printf("Logging into: %s\n", os.Getenv("REDDIT_USERNAME"))
	cookie, modhash, err := MyLoginSession(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		"Resistance Telegram Botter",
	)
	if err != nil {
		errorLogger("Error logging into reddit! " + err.Error())
	} else {
		log.Println(fmt.Sprintf("Succesfully logged in."))
	}

	return cookie, modhash
}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}
	log.Printf("Starting bot using port %s\n", port)

	errorMessages := []string{}
	var errorLogger = func(msg string) {
		log.Println(msg)
		newMsg := [...]string{getTime() + ": " + msg}
		errorMessages = append(newMsg[:], errorMessages...)
	}

	teleBot := Telegram{
		key:         os.Getenv("TELE_KEY"),
		errorLogger: errorLogger,
		lastUpdate:  0,
	}

	go listenForUpdates(teleBot, errorLogger)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, &errorMessages)
	r.Run(":" + port)

}
