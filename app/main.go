// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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

func trapCommand(url string, update Update) {
	if strings.ContainsAny(update.Message.Text, "traps") {
		go sendMessage("https://www.youtube.com/watch?v=9E1YYSZ9qrk", url, update)
	}
}

func executeCommand(url string, update Update) {
	trapCommand(url, update)
}

// Create our routes
func initRoutes(router *gin.Engine, teleurl string, logger func(*gin.Context)) {

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Telegram Bot")
	})

	router.POST("/new-message", func(c *gin.Context) {

		var update Update
		bindError := c.BindJSON(&update)

		if bindError == nil {
			c.JSON(http.StatusOK, update)
			executeCommand(teleurl, update)
		} else {
			log.Println(bindError)
			c.JSON(http.StatusBadRequest, update)
		}

	})

}

func main() {

	port := os.Getenv("PORT")

	// Can't run a server without a port
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

	// Set Logger
	var logger func(*gin.Context)
	if os.Getenv("APP_ENV") == "production" {
		log.Println("Running api server in production mode")
		logger = func(c *gin.Context) {}
	} else {
		log.Println("Running api server in debug mode")
		logger = func(c *gin.Context) {
			var data, err = c.GetRawData()
			if err == nil {
				fmt.Println("\nData recieved: ")
				fmt.Println(string(data))
			} else {
				fmt.Println("Error attempting to log data")
			}
		}
	}

	// Start server
	teleurl := "https://api.telegram.org/bot" + os.Getenv("TELE_KEY") + "/"
	initRoutes(r, teleurl, logger)
	r.Run(":" + port)

}
