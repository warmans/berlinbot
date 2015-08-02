package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/google/go-querystring/query"
	"github.com/warmans/berlinbot/songkick"
	"text/template"
)

const(
	SUBMIT_TYPE_SELF = "self"
	SUBMIT_TYPE_LINK = "link"
)

type Bot struct {
	Reddit *RedditClient
	Songkick *songkick.SongkickClient
}

func (b *Bot) GetMe() {

	log.Print(`Fetching me...`)

	request, err := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := b.Reddit.DoAuthorized(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Print(string(bodyBytes))
}

func (b *Bot) SubmitEvents() {

	tpl, err := template.ParseFiles(`post.tpl.md`)
	if err != nil {
		log.Printf(`Failed to render template: %s`, err.Error())
		return
	}

	out := bytes.NewBufferString(``)
	upcoming := b.Songkick.GetUpcoming(28443, 1)

	if err := tpl.Execute(out, upcoming); err != nil {
		log.Printf(`failed to execute template: %s`, err.Error())
	}

	b.Submit(SUBMIT_TYPE_SELF, `Events Monthly`, out.String(), `test`)
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

	response, err := b.Reddit.DoAuthorized(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response)

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Print(string(bodyBytes))
}
