// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Resposible for sending a message to the appropriate group chat
func sendMessage(message string, url string, update Update) {

	// Send Message to telegram's api
	resp, err := http.Post(url+"sendMessage", "application/json", bytes.NewBuffer([]byte(`{
		"chat_id": `+strconv.FormatInt(update.Message.Chat.ID, 10)+`,
		"text": "`+message+`"
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
func getCommands(url string) []func(Update) {

	killStatements := []string{
		" teabagged a piranha tank",
		" died of a heart attack while watching hentia",
		" just got back from yiffing",
		" was bitten by a horse",
		" was bent over and given a slap on the bottom",
		" changed their major to BIS",
		" got drilled",
		" got paddeled",
		" drank bleach",
		" was put on a group project with kleiderar",
		" got hired by lockheed",
		" was forever shunned by the resume gods",
	}

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
	}
}

// Create our routes
func initRoutes(router *gin.Engine, teleurl string) {

	commands := getCommands(teleurl)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Telegram Bot is running!")
	})

	router.POST("/new-message", func(c *gin.Context) {

		var update Update
		bindError := c.BindJSON(&update)

		if bindError == nil {
			c.JSON(http.StatusOK, update)

			// Executes commands
			for _, command := range commands {
				command(update)
			}
		} else {
			log.Println(bindError)
			c.JSON(http.StatusBadRequest, update)
		}

	})

}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}

	log.Println("Starting bot...")

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	// Start server
	teleurl := "https://api.telegram.org/bot" + os.Getenv("TELE_KEY") + "/"
	initRoutes(r, teleurl)
	r.Run(":" + port)

}
