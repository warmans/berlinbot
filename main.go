package main

import (
	"flag"
	"net/http"
)

func main() {

	var appId string
	var appSecret string
	var username string
	var password string

	flag.StringVar(&appId, `id`, ``, `app id`)
	flag.StringVar(&appSecret, `sec`, ``, `app secret key`)
	flag.StringVar(&username, `user`, ``, `user to post as`)
	flag.StringVar(&password, `pass`, ``, `user password`)

	client := &Client{HTTPClient: &http.Client{}}
	client.CreateToken(appId, appSecret, username, password)

	bot := &Bot{Client: client}
	bot.Submit(SUBMIT_TYPE_SELF, `testing...`, `testing...`, `test`)
}
