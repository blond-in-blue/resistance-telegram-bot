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

var rule34Command = BotCommand{
	Name: "Rule34",
	Description: "Search reddit's rule 34: /rule34 golang",
	Matcher: messageContainsCommandMatcher("rule34"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		searchTerms := getContentFromCommand(update.Message.Text, "rule34")

		if searchTerms != "" {
			submissions, err := bot.redditUser.SearchSubreddit("rule34", searchTerms)

			if err != nil {
				bot.errorReport.Log("Error searching subreddit: " + err.Error())
				respChan <- *NewTextBotResponse("Error searching subreddit", update.Message.Chat.ID)
			}

			if len(submissions) > 0 {
				respChan <- *NewTextBotResponse(submissions[0].Title+"\n"+submissions[0].URL, update.Message.Chat.ID)
			} else {
				respChan <- *NewTextBotResponse(fmt.Sprintf("Didn't find anything for '%s'", searchTerms), update.Message.Chat.ID)
			}
		}
	},
}


var hedgehogCommand = BotCommand{
	Name: "Hedgehog",
	Description: "Have you been hedgehogged? /hedgehog eli",
	Matcher: messageContainsCommandMatcher("hedgehog"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		searchTerms := getContentFromCommand(update.Message.Text, "hedgehog")

		if searchTerms != "" {

			submissions, err := bot.redditUser.SearchSubreddit("thehedgehog", searchTerms)

			if err != nil {
				bot.errorReport.Log("Error searching subreddit: " + err.Error())
				respChan <- *NewTextBotResponse("Error searching subreddit", update.Message.Chat.ID)
			}

			if len(submissions) > 0 {
				respChan <- *NewTextBotResponse(submissions[0].Title+"\n"+submissions[0].URL, update.Message.Chat.ID)
			} else {
				respChan <- *NewTextBotResponse(fmt.Sprintf("'%s' has not been hedgehogged ", searchTerms), update.Message.Chat.ID)
			}
		}
	},
}

var saveCommand = BotCommand{
	Name: "Save",
	Description: "Save a text post to the subreddit",
	Matcher: messageContainsCommandMatcher("save"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		term := getContentFromCommand(update.Message.Text, "save")

		if term == "" {
			respChan <- *NewTextBotResponse("Please provide a title for the post", update.Message.Chat.ID)
			return
		}

		if update.Message.ReplyToMessage == nil {
			respChan <- *NewTextBotResponse("Reply to a message and say save to save to the subreddit", update.Message.Chat.ID)
			return
		}
	
		if update.Message.ReplyToMessage.Text == "" {
			respChan <- *NewTextBotResponse("I can only save text, give me some text or open up a feature branch", update.Message.Chat.ID)
			return
		}
	
		log.Printf("update: %s", update.Message.ReplyToMessage.Text)
	
		info, err := bot.redditUser.PostToSubreddit(fmt.Sprintf("%s:\n\n%s", update.Message.ReplyToMessage.From.UserName, update.Message.ReplyToMessage.Text), term, "smartestretards")
		if err != nil {
			bot.errorReport.Log("Unable to post to reddit: " + err.Error())
			respChan <- *NewTextBotResponse("Unable to post to reddit", update.Message.Chat.ID)
		} else {
			respChan <- *NewTextBotResponse(info, update.Message.Chat.ID)
		}
	},
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
