package main

import (
	"flag"
	"net/http"
	"github.com/warmans/berlinbot/songkick"
	"log"
	"github.com/warmans/berlinbot/spotify"
	"github.com/warmans/berlinbot/bot"
)

func main() {

	var dbg bool

	var redditAppId string
	var redditAppSecret string
	var redditUsername string
	var redditPassword string
	var redditSub string
	var redditTitle string

	var songkickKey string

	flag.BoolVar(&dbg, `dbg`, false, `Dry run (dump to stdout instead of posting)`)

	//reddit details required
	flag.StringVar(&redditAppId, `reddit-id`, ``, `reddit app id`)
	flag.StringVar(&redditAppSecret, `reddit-key`, ``, `reddit app secret key`)
	flag.StringVar(&redditUsername, `reddit-user`, ``, `reddit user to post as`)
	flag.StringVar(&redditPassword, `reddit-pass`, ``, `reddit user password`)
	flag.StringVar(&redditSub, `reddit-sub`, `test`, `subreddit to post to`)
	flag.StringVar(&redditTitle, `reddit-title`, `Events This Week`, `Title of post`)

	//songkick details required
	flag.StringVar(&songkickKey, `songkick-key`, ``, `songkick api key`)
	flag.Parse()

	if redditAppId == "" || redditAppSecret == "" || redditUsername == "" || redditPassword == "" {
		log.Fatal("All reddit args must be supplied")
	}

	if songkickKey == "" {
		log.Fatal("Songkick key must be provided")
	}

	songkick := &songkick.SongkickClient{APIKey: songkickKey, HTTPClient: &http.Client{}}
	spotify := &spotify.Spotify{&http.Client{}}


	reddit := &bot.RedditClient{HTTPClient: &http.Client{}}
	if dbg == false {
		//strict API limits for this. Only do it when it's used
		reddit.CreateToken(redditAppId, redditAppSecret, redditUsername, redditPassword)
	}

	//execute
	bot := &bot.Bot{Debug: dbg, Reddit: reddit, Songkick: songkick, Spotify: spotify}
	bot.SubmitEvents(redditSub, redditTitle)
}
