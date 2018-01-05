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
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// Telegram struct for taking care of commands
type Telegram struct {
	key         string
	errorLogger func(string)
	lastUpdate  int
}

func (telebot Telegram) getURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/", telebot.key)
}

// SendMessage Resposible for sending a message to the appropriate group chat
func (telebot Telegram) SendMessage(message string, chatID int64) {

	// Send Message to telegram's api
	resp, err := http.Post(telebot.getURL()+"sendMessage", "application/json", bytes.NewBuffer([]byte(`{
		"chat_id": `+strconv.FormatInt(chatID, 10)+`,
		"text": "`+message+`",
		"parse_mode": "HTML"
	}`)))

	// Catch errors
	if err != nil {
		telebot.errorLogger("Error sending message: " + err.Error())
		return
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		telebot.errorLogger("Error recieving response from tele: " + err.Error())
	}
}

// GetUpdates queries telegram for latest updates
func (telebot *Telegram) GetUpdates() ([]Update, error) {
	resp, err := http.Get(telebot.getURL() + "getUpdates?offset=" + strconv.Itoa(telebot.lastUpdate))

	// Sometimes Telegram will just randomly send a 502
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	defer resp.Body.Close()

	var updates BatchUpdates
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &updates)
	if err != nil {
		return nil, err
	}

	for _, update := range updates.Result {
		telebot.lastUpdate = update.UpdateID + 1
	}

	return updates.Result, nil
}

func (telebot Telegram) sendImage(path string, chatID int64) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer

	w.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	file, err := os.Open(path)
	if err != nil {
		telebot.errorLogger(err.Error())
	}

	img, _, err := image.Decode(file)
	if err != nil {
		telebot.errorLogger(err.Error())
	}

	if fw, err = w.CreateFormFile("photo", "image.png"); err != nil {
		telebot.errorLogger(err.Error())
	}
	if err = png.Encode(fw, img); err != nil {
		telebot.errorLogger(err.Error())
	}

	w.CreateFormField("something")

	w.Close()

	req, err := http.NewRequest("POST", telebot.getURL()+"sendPhoto", &b)
	if err != nil {
		telebot.errorLogger(err.Error())
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		telebot.errorLogger(err.Error())
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		telebot.errorLogger(err.Error())
	}
	log.Println(string(bytes))
}
