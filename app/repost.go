package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func PersonWhoPostedFile(chatId string, fileId string) (string, error) {
	f, err := os.OpenFile("repost/"+chatId, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		if fileId == entry[0] {
			return entry[2], nil
		}
	}

	return "", nil
}

func PersonWhoPostedHash(chatId string, hash string) (string, error) {
	f, err := os.OpenFile("repost/"+chatId, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		if hash == entry[1] {
			return entry[2], nil
		}
	}

	return "", nil
}

func HashFile(filePath string) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func StoreFileEntry(chatId string, poster string, hash string, fileId string) error {
	f, err := os.OpenFile("repost/"+chatId, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	str := fmt.Sprintf("%s %s %s\n", fileId, hash, poster)
	log.Println(str)
	if _, err = f.WriteString(str); err != nil {
		return err
	}

	return nil
}

var repostCommand = BotCommand{
	Name:        "Repost",
	Description: "Attempts to detect reposted images",
	Matcher: func(update Update) bool {
		return update.Message.Photo != nil
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		photos := *update.Message.Photo
		chadId := strconv.FormatInt(update.Message.Chat.ID, 10)
		poster, err := PersonWhoPostedFile(chadId, photos[0].FileID)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if poster != "" {
			respChan <- *NewTextBotResponse(fmt.Sprintf("REPOST: %s has already posted this", poster), update.Message.Chat.ID)
		} else {
			bot.GetFile(photos[0].FileID)
			hash, err := HashFile(fmt.Sprintf("media/%s", photos[0].FileID))
			if err == nil {
				hashPoster, err := PersonWhoPostedHash(chadId, hash)
				if err != nil {
					log.Println(err.Error())
					return
				}
				if hashPoster != "" {
					respChan <- *NewTextBotResponse(fmt.Sprintf("REPOST: %s has already posted this", hashPoster), update.Message.Chat.ID)
				} else {
					err := StoreFileEntry(chadId, update.Message.From.UserName, hash, photos[0].FileID)
					if err != nil {
						log.Println(err.Error())
					}
				}
			} else {
				log.Println(err.Error())
			}
		}
	},
}
