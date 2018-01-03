// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
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

func sendImage(url string, update Update, errorLogger func(string)) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer

	w.WriteField("chat_id", strconv.FormatInt(update.Message.Chat.ID, 10))
	file, err := os.Open("media/out.png")
	if err != nil {
		errorLogger(err.Error())
	}

	img, msg, err := image.Decode(file)
	if err != nil {
		errorLogger(err.Error())
	}
	log.Println(msg)

	if fw, err = w.CreateFormFile("photo", "image.png"); err != nil {
		errorLogger(err.Error())
	}
	if err = png.Encode(fw, img); err != nil {
		errorLogger(err.Error())
	}

	w.CreateFormField("something")

	w.Close()

	req, err := http.NewRequest("POST", url+"sendPhoto", &b)
	if err != nil {
		errorLogger(err.Error())
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errorLogger(err.Error())
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		errorLogger(err.Error())
	}
	log.Println(string(bytes))
}

// Resposible for sending a message to the appropriate group chat
func sendMessage(message string, url string, update Update) {

	// Send Message to telegram's api
	resp, err := http.Post(url+"sendMessage", "application/json", bytes.NewBuffer([]byte(`{
		"chat_id": `+strconv.FormatInt(update.Message.Chat.ID, 10)+`,
		"text": "`+message+`",
		"parse_mode": "HTML"
	}`)))

	// Catch errors
	if err != nil {
		log.Println("Error sending message:")
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	// Read and print message
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("\nTelegram Said: ")
	log.Println(string(body))
}

// Builds and returns commands with url.
func getCommands(url string, errorLogger func(string), redditSession *http.Cookie, modhash string) []func(Update) {

	return []func(update Update){

		// Kill command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "kill")
			if len(commands) > 1 {
				n := rand.Int() % len(killStatements)
				go sendMessage(strings.TrimSpace(commands[1])+killStatements[n], url, update)
			}
		},

		// Traps command
		func(update Update) {
			if strings.Contains(update.Message.Text, "traps") {
				go sendMessage("https://www.youtube.com/watch?v=9E1YYSZ9qrk", url, update)
			}
		},

		// God command
		func(update Update) {
			if strings.Contains(update.Message.Text, "gg") {
				go sendMessage("GOD IS GREAT", url, update)
			}
		},

		// Rule34 command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "rule34")
			if len(commands) > 1 {
				go rule34Search(strings.TrimSpace(commands[1]), url, update, errorLogger, redditSession)
			}
		},

		// Save command
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "save")
			if len(commands) > 1 {
				go SaveCommand(strings.TrimSpace(commands[1]), url, update, errorLogger, redditSession, modhash)
			}
		},

		//pokedexSerach
		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "pokedex")
			if len(commands) > 1 {
				go pokedexSerach(strings.TrimSpace(commands[1]), url, update, errorLogger)
			}
		},

		func(update Update) {
			commands := strings.SplitAfter(update.Message.Text, "murder")
			if len(commands) > 1 {
				log.Println()
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
				text := strings.TrimSpace(commands[1])
				dc.DrawStringAnchored(text, 500, 120, 0.0, 0.0)
				dc.SavePNG("media/out.png")
				sendImage(url, update, errorLogger)
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

func listenForUpdates(teleurl string, errorLogger func(string)) {
	var lastUpdate = -1

	cookie, modhash := logginToReddit(errorLogger)
	commands := getCommands(teleurl, errorLogger, cookie, modhash)

	for {
		// Sleep first, so if we error out and continue to the next loop, we still end up waiting
		time.Sleep(time.Second)

		resp, err := http.Get(teleurl + "getUpdates?offset=" + strconv.Itoa(lastUpdate))

		// Sometimes Telegram will just randomly send a 502
		if err != nil || resp.StatusCode != 200 {
			errorLogger("Error Obtaining Updates: " + err.Error())
			continue
		}

		defer resp.Body.Close()

		var updates BatchUpdates
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorLogger("Error Reading Body: " + err.Error())
			continue
		}

		err = json.Unmarshal([]byte(body), &updates)
		if err != nil {
			errorLogger("Error Parsing Telegram getUpdates Response: " + err.Error() + "; Response body: " + string(body))
			continue
		}

		// Dispatch incoming messages to appropriate functions
		for _, update := range updates.Result {
			if update.Message != nil {
				log.Println("Msg: " + update.Message.Text)
				for _, command := range commands {
					command(update)
				}
			}
			lastUpdate = update.UpdateID + 1
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

	teleurl := "https://api.telegram.org/bot" + os.Getenv("TELE_KEY") + "/"

	go listenForUpdates(teleurl, errorLogger)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r, &errorMessages)
	r.Run(":" + port)

}
