package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

// NewTelegramBot Creates a new telegram bot
func NewTelegramBot(key string, errorLogger func(string)) *Telegram {
	t := Telegram{
		key:         key,
		errorLogger: errorLogger,
		lastUpdate:  0,
		url:         fmt.Sprintf("https://api.telegram.org/bot%s/", key),
	}
	return &t
}

// Telegram struct for taking care of commands
type Telegram struct {
	key         string
	errorLogger func(string)
	lastUpdate  int
	url         string
}

// SendMessage Resposible for sending a message to the appropriate group chat
func (telebot Telegram) SendMessage(message string, chatID int64) {

	// Send Message to telegram's api
	resp, err := http.Post(telebot.url+"sendMessage", "application/json", bytes.NewBuffer([]byte(`{
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
	resp, err := http.Get(telebot.url + "getUpdates?offset=" + strconv.Itoa(telebot.lastUpdate))

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

// GetImage Downloads an image by its file id and returns the filepath on the system
func (telebot Telegram) GetImage(fileID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%sgetFile?file_id=%s", telebot.url, fileID))

	log.Println("Begining download")

	if err != nil {
		log.Println("Error: " + err.Error())
		return "", err
	}

	var imageResponse GetImageResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(body), &imageResponse)
	if err != nil {
		log.Println("err: " + err.Error())
		return "", err
	}

	if imageResponse.Ok == false {
		return "", errors.New("telegram resolved unsucessfully")
	}

	resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", telebot.key, imageResponse.Result.FilePath))

	if err != nil {
		return "", err
	}

	filePath := "media/" + fileID

	output, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(output, resp.Body)

	if err != nil {
		return "", err
	}

	log.Println("Succesfully downloaded")

	return filePath, nil
}

func (telebot Telegram) deleteMessage(chatID int64, messageID int) (bool, error) {

	resp, err := http.Get(fmt.Sprintf("%sdeleteMessage?chat_id=%s&message_id=%d", telebot.url, strconv.FormatInt(chatID, 10), messageID))

	if err != nil {
		log.Println("Error: " + err.Error())
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf(string(body))
	return true, nil

}

// SendPhotoByID send photo by already existing file id
func (telebot Telegram) SendPhotoByID(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendPhoto?chat_id=%s&photo=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot Telegram) SendSticker(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendSticker?chat_id=%s&sticker=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
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

	req, err := http.NewRequest("POST", telebot.url+"sendPhoto", &b)
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
