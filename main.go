package main

import (
	"flag"
	"net/http"
	"github.com/warmans/berlinbot/songkick"
)

func main() {

	var redditAppId string
	var redditAppSecret string
	var redditUsername string
	var redditpassword string

	var songkickKey string

	//reddit details required
	flag.StringVar(&redditAppId, `reddit-id`, ``, `reddit app id`)
	flag.StringVar(&redditAppSecret, `reddit-key`, ``, `reddit app secret key`)
	flag.StringVar(&redditUsername, `reddit-user`, ``, `reddit user to post as`)
	flag.StringVar(&redditpassword, `reddit-pass`, ``, `reddit user password`)

	//songkick details required
	flag.StringVar(&songkickKey, `songkick-key`, ``, `songkick api key`)

	flag.Parse()

	reddit := &RedditClient{HTTPClient: &http.Client{}}
	reddit.CreateToken(redditAppId, redditAppSecret, redditUsername, redditpassword)

	songkick := &songkick.SongkickClient{APIKey: songkickKey, HTTPClient: &http.Client{}}

	//execute
	bot := &Bot{Reddit: reddit, Songkick: songkick}
	bot.SubmitEvents()
}
