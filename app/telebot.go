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

// TeleBot talks to telegram and manages application state
type TeleBot struct {
	key         string
	lastUpdate  int
	url         string
	chatBuffers map[string]MessageStack
	chatAliases map[string]string
	errorReport Report
	redditUser  RedditAccount
}

// NewTelegramBot Creates a new telegram bot
func NewTelegramBot(key string, errorReport Report, redditAccount RedditAccount) *TeleBot {
	t := TeleBot{
		key:         key,
		lastUpdate:  0,
		url:         fmt.Sprintf("https://api.telegram.org/bot%s/", key),
		errorReport: errorReport,
		chatBuffers: make(map[string]MessageStack),
		chatAliases: make(map[string]string),
		redditUser:  redditAccount,
	}
	return &t
}

func (telebot TeleBot) IsAliasSet(alias string) (string, bool) {
	str, b := telebot.chatAliases[alias]
	return str, b
}

func (telebot TeleBot) SetChatAlias(alias string, chatID int64) {
	telebot.chatAliases[alias] = strconv.FormatInt(chatID, 10)
}

// PushMessageToChatBuffer moves a message to the appropriate chats
func (telebot *TeleBot) PushMessageToChatBuffer(lookup string, message Message) {
	location := strconv.FormatInt(message.Chat.ID, 10)
	if lookup != "" {
		location = lookup
		alias, exists := telebot.chatAliases[location]
		if exists {
			location = alias
		}
	}
	telebot.chatBuffers[location] = telebot.chatBuffers[location].Push(message)
}

func (telebot *TeleBot) ClearBuffer(chatID int64) <-chan Message {
	lookup := strconv.FormatInt(chatID, 10)
	buffer := telebot.chatBuffers[lookup]
	telebot.chatBuffers[lookup] = make([]Message, 0)
	return buffer.Everything()
}

// ChatBuffer returns the buffer for a specific chat given the lookup
func (telebot TeleBot) ChatBuffer(lookup string) MessageStack {

	// Try looking it up immediately like they gave us a chat id
	buffer, exists := telebot.chatBuffers[lookup]
	if exists {
		return buffer
	}

	// If we didn't find anything, they might of given us an alias
	alias, exists := telebot.chatAliases[lookup]
	if exists {
		return telebot.chatBuffers[alias]
	}

	return nil
}

// SendMessage Resposible for sending a message to the appropriate group chat
func (telebot TeleBot) SendMessage(message string, chatID int64) {

	postValues := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"text":       message,
		"parse_mode": "HTML",
	}

	jsonValue, err := json.Marshal(postValues)

	if err != nil {
		telebot.errorReport.Log("Error encoding json: " + err.Error())
		return
	}

	req, err := http.NewRequest("POST", telebot.url+"sendMessage", bytes.NewBuffer(jsonValue))

	if err != nil {
		telebot.errorReport.Log("Error creating message: " + err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		telebot.errorReport.Log("Error sending message: " + err.Error())
		return
	}

	// Catch errors
	if err != nil {
		telebot.errorReport.Log("Error sending message: " + err.Error())
		return
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		telebot.errorReport.Log("Error recieving response from tele: " + err.Error())
	}
}

// GetUpdates queries telegram for latest updates
func (telebot *TeleBot) GetUpdates() ([]Update, error) {
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
func (telebot TeleBot) GetImage(fileID string) (string, error) {
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

func (telebot TeleBot) deleteMessage(chatID int64, messageID int) (bool, error) {

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
func (telebot TeleBot) SendPhotoByID(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendPhoto?chat_id=%s&photo=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) SendSticker(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendSticker?chat_id=%s&sticker=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) sendImage(path string, chatID int64) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer

	w.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	file, err := os.Open(path)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	img, _, err := image.Decode(file)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	if fw, err = w.CreateFormFile("photo", "image.png"); err != nil {
		telebot.errorReport.Log(err.Error())
	}
	if err = png.Encode(fw, img); err != nil {
		telebot.errorReport.Log(err.Error())
	}

	w.CreateFormField("something")

	w.Close()

	req, err := http.NewRequest("POST", telebot.url+"sendPhoto", &b)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		telebot.errorReport.Log(err.Error())
	}
	log.Println(string(bytes))
}
