package main

import (
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"time"
	"net/url"
)

type RedditClient struct {
	HTTPClient *http.Client
	AuthResponse *AuthResponse
}

func (c *RedditClient) CreateToken(appId, appSecret, user, password string) {

	bodyBuff := bytes.NewBufferString(`grant_type=password&username=`+url.QueryEscape(user)+`&password=`+url.QueryEscape(password)+``)
	request, err := http.NewRequest(`POST`, `https://www.reddit.com/api/v1/access_token`, bodyBuff)
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(appId, appSecret)
	request.Header.Set(`Content-type`, `application/x-www-form-urlencoded`)

	response, err := c.HTTPClient.Do(request)
	defer response.Body.Close()

	authResponse := &AuthResponse{Created:int(time.Now().Unix())}
	if err := json.NewDecoder(response.Body).Decode(authResponse); err != nil {
		log.Fatal(err)
	}

	c.AuthResponse = authResponse
}

func (c *RedditClient) DoAuthorized(request *http.Request) (*http.Response, error) {
	if c.AuthResponse.AccessToken == "" {
		log.Fatal(`Create a token first`)
	}

	log.Print(`append token: `+c.AuthResponse.TokenType+` `+c.AuthResponse.AccessToken)
	request.Header.Set(`Authorization`, c.AuthResponse.TokenType+` `+c.AuthResponse.AccessToken)
	request.Header.Set(`User-Agent`, `BerlinBot/0.1 by warmans`)

	return c.HTTPClient.Do(request)
}

type AuthResponse struct {
	Created     int
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type SubmitRequest struct {
	Title    string `url:"title"`
	Text     string `url:"text"`
	SR       string `url:"sr"`
	Kind     string `url:"kind"`
	Resubmit bool   `url:"resubmit"`
}
