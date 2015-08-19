package spotify

import (
	"net/http"
	"net/url"
	"log"
	"encoding/json"
)

type Artist struct {
	ExternalUrls struct {
					 Spotify string `json:"spotify"`
				 } `json:"external_urls"`
	Followers    struct {
					 Href  interface{} `json:"href"`
					 Total int `json:"total"`
				 } `json:"followers"`
	Genres       []string `json:"genres"`
	Href         string `json:"href"`
	ID           string `json:"id"`
	Images       []struct {
		Height int `json:"height"`
		URL    string `json:"url"`
		Width  int `json:"width"`
	} `json:"images"`
	Name         string `json:"name"`
	Popularity   int `json:"popularity"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type Result struct {
	Artists struct {
				Href     string `json:"href"`
				Items    []Artist `json:"items"`
				Limit    int `json:"limit"`
				Next     interface{} `json:"next"`
				Offset   int `json:"offset"`
				Previous interface{} `json:"previous"`
				Total    int `json:"total"`
			} `json:"artists"`
}

type Spotify struct {
	HTTPClient *http.Client
}

func (s Spotify) GuessArtist(name string) *Artist {
	data, err := s.HTTPClient.Get(`https://api.spotify.com/v1/search?q=` + url.QueryEscape(name) + `&type=artist`)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	defer data.Body.Close()

	var result Result
	json.NewDecoder(data.Body).Decode(&result)

	if len(result.Artists.Items) == 0 {
		return nil
	}

	return &result.Artists.Items[0]
}

