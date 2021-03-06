package bot

import (
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"time"
	"net/url"
	"github.com/google/go-querystring/query"
	"io/ioutil"
)

const(
	SUBMIT_TYPE_SELF = "self"
	SUBMIT_TYPE_LINK = "link"
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

func (r *RedditClient) Submit(kind, title, text, subreddit string) {
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

	response, err := r.DoAuthorized(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(response)

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Print(string(bodyBytes))
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
