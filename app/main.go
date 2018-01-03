// Structs from:
// https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

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

func rule34Search(term string, url string, update Update, errorLogger func(string), redditSession *http.Cookie) {
	log.Println("searching rule 34: " + term)

	// Create a request to be sent to reddit
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/rule34/search.json?q="+term+"&restrict_sr=on&sort=relevance&t=all", nil)
	if err != nil {
		errorLogger("Error creating reddit request: " + err.Error())
	}
	req.Header.Set("User-Agent", "Resistance Telegram Bot")
	req.AddCookie(redditSession)

	// Actually query reddit
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errorLogger("Error querying reddit: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorLogger("Error querying reddit: " + resp.Status)
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorLogger("Error querying reddit: " + err.Error())
	}

	log.Println(string(respbytes))

	body := bytes.NewBuffer(respbytes)

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)

	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		errorLogger(err.Error())
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	if len(submissions) > 0 {
		sendMessage(submissions[0].URL, url, update)
	}

}

// Builds and returns commands with url.
func getCommands(url string, errorLogger func(string), redditSession *http.Cookie) []func(Update) {

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
				dc := gg.NewContext(1000, 1000)
				dc.DrawCircle(500, 500, 400)
				dc.SetRGB(0, 0, 0)
				dc.Fill()
				dc.SavePNG("out.png")
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

}

func listenForUpdates(teleurl string, errorLogger func(string)) {
	var lastUpdate = -1

	commands := getCommands(teleurl, errorLogger, logginToReddit(errorLogger))

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

// MyLoginSession creates a new session for those who want to log into a
func MyLoginSession(username, password, useragent string) (*http.Cookie, string, error) {

	loginURL := fmt.Sprintf("https://www.reddit.com/api/login/%s", username)
	postValues := url.Values{
		"user":     {username},
		"passwd":   {password},
		"api_type": {"json"},
	}

	// Build our request
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", useragent)

	http.DefaultClient.Timeout = time.Second * 10
	log.Println(http.DefaultClient)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", err
	}

	// Get the session cookie.
	var redditCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "reddit_session" {
			redditCookie = cookie
		}
	}

	// Get the modhash from the JSON.
	type Response struct {
		JSON struct {
			Errors [][]string
			Data   struct {
				Modhash string
			}
		}
	}

	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return redditCookie, "", err
	}

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return redditCookie, "", errors.New(strings.Join(msg, ", "))
	}
	modhash := r.JSON.Data.Modhash

	return redditCookie, modhash, nil
}

func logginToReddit(errorLogger func(string)) *http.Cookie {

	log.Printf("Logging into: %s\n", os.Getenv("REDDIT_USERNAME"))
	cookie, modhash, err := MyLoginSession(
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		"Resistance Telegram Botter",
	)
	if err != nil {
		errorLogger("Error logging into reddit! " + err.Error())
	} else {
		log.Println(fmt.Sprintf("Succesfully logged in %s", modhash))
	}

	return cookie
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

	log.Println("Get Ready....")
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
