package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RedditAccount struct {
	cookie  *http.Cookie
	modhash string
}

func LoginToReddit(username, password, useragent string) (RedditAccount, error) {
	loginURL := fmt.Sprintf("https://www.reddit.com/api/login/%s", username)
	postValues := url.Values{
		"user":     {username},
		"passwd":   {password},
		"api_type": {"json"},
	}

	// Build our request
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return RedditAccount{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", useragent)

	http.DefaultClient.Timeout = time.Second * 10
	log.Println(http.DefaultClient)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RedditAccount{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RedditAccount{}, err
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
		return RedditAccount{}, err
	}

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return RedditAccount{}, errors.New(strings.Join(msg, ", "))
	}
	modhash := r.JSON.Data.Modhash

	return RedditAccount{redditCookie, modhash}, nil
}

func rule34Search(term string, telebot TeleBot, update Update) {

	submissions, err := telebot.redditUser.SearchSubreddit("rule34", term)

	if err != nil {
		telebot.errorReport.Log("Error searching subreddit: " + err.Error())
		telebot.SendMessage("Error searching subreddit", update.Message.Chat.ID)
		return
	}

	if len(submissions) > 0 {
		telebot.SendMessage(submissions[0].Title+"\n"+submissions[0].URL, update.Message.Chat.ID)
	} else {
		telebot.SendMessage(fmt.Sprintf("Didn't find anything for '%s'", term), update.Message.Chat.ID)
	}

}

func hedgeHogCommand(term string, telebot TeleBot, update Update) {

	submissions, err := telebot.redditUser.SearchSubreddit("thehedgehog", term)

	if err != nil {
		telebot.errorReport.Log("Error searching subreddit: " + err.Error())
		telebot.SendMessage("Error searching", update.Message.Chat.ID)
		return
	}

	if len(submissions) > 0 {
		telebot.SendMessage(submissions[0].URL, update.Message.Chat.ID)
	} else {
		telebot.SendMessage(fmt.Sprintf("'%s' has not been hedgehogged ", term), update.Message.Chat.ID)
	}

}

// SaveCommand posts to our subreddit
func SaveCommand(term string, telebot TeleBot, update Update) {

	if update.Message.ReplyToMessage == nil {
		telebot.SendMessage("Reply to a message and say save to save to the subreddit", update.Message.Chat.ID)
		return
	}

	if update.Message.ReplyToMessage.Text == "" {
		telebot.SendMessage("I can only save text, give me some text or open up a feature branch", update.Message.Chat.ID)
		return
	}

	//log.Println("Going to save... " + term)
	log.Printf("update: %s", update.Message.ReplyToMessage.Text)

	info, err := telebot.redditUser.PostToSubreddit(fmt.Sprintf("%s:\n\n%s", update.Message.ReplyToMessage.From.UserName, update.Message.ReplyToMessage.Text), term, "smartestretards")
	if err != nil {
		telebot.errorReport.Log("Unable to post to reddit: " + err.Error())
		telebot.SendMessage("Unable to post to reddit", update.Message.Chat.ID)
	} else {
		telebot.SendMessage(info, update.Message.Chat.ID)
	}
}

func (r RedditAccount) PostToSubreddit(textPost string, title string, subreddit string) (string, error) {
	// Create a request to be sent to reddit
	vals := &url.Values{
		"title":       {title},
		"url":         {textPost},
		"text":        {textPost},
		"sr":          {subreddit},
		"kind":        {"self"},
		"sendreplies": {"true"},
		"resubmit":    {"true"},
		"extension":   {"json"},
		"uh":          {r.modhash},
	}
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/submit?"+vals.Encode(), nil)
	if err != nil {
		//errorLogger("Error creating reddit request: " + err.Error())
		//telebot.SendMessage("Error: "+err.Error(), update.Message.Chat.ID)
		return "", err
	}
	req.Header.Set("User-Agent", "Resistance Telegram Bot")
	req.AddCookie(r.cookie)

	// Actually query reddit
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body := bytes.NewBuffer(respbytes)
	log.Println(body)
	if strings.Contains(string(respbytes), "error") {
		return "", err
	} else {
		return "Success", nil
	}
}

func (reddit RedditAccount) SearchSubreddit(subreddit string, term string) ([]*Submission, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.reddit.com/r/%s/search.json?q=%s&restrict_sr=on&include_over_18=on&sort=relevance&t=all", subreddit, strings.Replace(term, " ", "+", -1)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Resistance Telegram Bot")
	req.AddCookie(reddit.cookie)

	// Actually query reddit
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}
	return submissions, nil
}
