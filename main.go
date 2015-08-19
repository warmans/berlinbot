package main

import (
	"flag"
	"net/http"
	"github.com/warmans/berlinbot/songkick"
	"log"
	"github.com/warmans/berlinbot/spotify"
	"fmt"
)

func main() {

	var redditAppId string
	var redditAppSecret string
	var redditUsername string
	var redditPassword string
	var redditSub string
	var redditTitle string

	var songkickKey string

	//reddit details required
	flag.StringVar(&redditAppId, `reddit-id`, ``, `reddit app id`)
	flag.StringVar(&redditAppSecret, `reddit-key`, ``, `reddit app secret key`)
	flag.StringVar(&redditUsername, `reddit-user`, ``, `reddit user to post as`)
	flag.StringVar(&redditPassword, `reddit-pass`, ``, `reddit user password`)
	flag.StringVar(&redditSub, `reddit-sub`, `test`, `subreddit to post to`)
	flag.StringVar(&redditTitle, `reddit-title`, `Concerts Weekly`, `Title of post`)

	//songkick details required
	flag.StringVar(&songkickKey, `songkick-key`, ``, `songkick api key`)
	flag.Parse()

	if redditAppId == "" || redditAppSecret == "" || redditUsername == "" || redditPassword == "" {
		log.Fatal("All reddit args must be supplied")
	}

	if songkickKey == "" {
		log.Fatal("Songkick key must be provided")
	}

	reddit := &RedditClient{HTTPClient: &http.Client{}}
	reddit.CreateToken(redditAppId, redditAppSecret, redditUsername, redditPassword)

	songkick := &songkick.SongkickClient{APIKey: songkickKey, HTTPClient: &http.Client{}}

	spotify := &spotify.Spotify{&http.Client{}}

	//execute
	bot := &Bot{Reddit: reddit, Songkick: songkick, Spotify: spotify}
	bot.SubmitEvents(redditSub, redditTitle)
}
