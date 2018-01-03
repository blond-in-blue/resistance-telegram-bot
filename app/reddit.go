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

func SaveCommand(term string, teleurl string, update Update, errorLogger func(string), redditSession *http.Cookie, modhash string) {

	if update.Message.ReplyToMessage == nil {
		sendMessage("Reply to a message and say save to save to the subreddit", teleurl, update)
		return
	}

	if update.Message.ReplyToMessage.Text == "" {
		sendMessage("I can only save text, give me some text or open up a feature branch", teleurl, update)
		return
	}

	log.Println("Going to save... " + term)
	log.Printf("update: %s", update.Message.ReplyToMessage.Text)

	// Create a request to be sent to reddit
	vals := &url.Values{
		"title":       {term},
		"url":         {update.Message.ReplyToMessage.Text},
		"text":        {update.Message.ReplyToMessage.Text},
		"sr":          {"smartestretards"},
		"kind":        {"self"},
		"sendreplies": {"true"},
		"resubmit":    {"true"},
		"extension":   {"json"},
		"uh":          {modhash},
	}
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/submit?"+vals.Encode(), nil)
	if err != nil {
		errorLogger("Error creating reddit request: " + err.Error())
		sendMessage("Error: "+err.Error(), teleurl, update)
	}
	req.Header.Set("User-Agent", "Resistance Telegram Bot")
	req.AddCookie(redditSession)

	// Actually query reddit
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errorLogger("Error querying reddit: " + err.Error())
		sendMessage("Error: "+err.Error(), teleurl, update)

	}
	defer resp.Body.Close()

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorLogger("Error reading body: " + err.Error())
		sendMessage("Error: "+err.Error(), teleurl, update)
	}

	// body := bytes.NewBuffer(respbytes)
	if strings.Contains(string(respbytes), "error") {
		errorLogger("failed to submit")
		sendMessage("failed to submit", teleurl, update)
	} else {
		sendMessage("I think it worked "+string(respbytes), teleurl, update)
	}
	// log.Println(body)
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
