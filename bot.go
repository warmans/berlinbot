package main

import (
	"log"
	"bytes"
	"github.com/warmans/berlinbot/songkick"
	"text/template"
	"github.com/warmans/berlinbot/spotify"
)

type Bot struct {
	Reddit   *RedditClient
	Songkick *songkick.SongkickClient
	Spotify  *spotify.Spotify
}

func (b *Bot) SubmitEvents(sub, title string) {

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

	b.Reddit.Submit(SUBMIT_TYPE_SELF, title, out.String(), sub)
}

