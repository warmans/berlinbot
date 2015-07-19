package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/google/go-querystring/query"
)

const(
	SUBMIT_TYPE_SELF = "self"
	SUBMIT_TYPE_LINK = "link"
)

type Bot struct {
	Client *Client
}

func (b *Bot) GetMe() {

	log.Print(`Fetching me...`)

	request, err := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := b.Client.DoAuthorized(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Print(string(bodyBytes))
}

func (b *Bot) Submit(kind, title, text, subreddit string) {
	log.Print(`Submitting...`)

	bodyContent := SubmitRequest{
		Kind: kind,
		Title: title,
		Text: text,
		SR: subreddit,
		Resubmit: true,
	}

	v, _ := query.Values(bodyContent)
	bodyBuffer := bytes.NewBufferString(v.Encode())

	request, err := http.NewRequest("POST", "https://oauth.reddit.com/api/submit", bodyBuffer)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set(`Content-type`, `application/x-www-form-urlencoded`)
	log.Print(request)

	response, err := b.Client.DoAuthorized(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response)

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Print(string(bodyBytes))
}
