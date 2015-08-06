package main

import (
	"log"
	"bytes"
	"github.com/warmans/berlinbot/songkick"
	"text/template"
)

type Bot struct {
	Reddit *RedditClient
	Songkick *songkick.SongkickClient
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

	b.Reddit.Submit(SUBMIT_TYPE_SELF, `Events Monthly`, out.String(), `test`)
}

