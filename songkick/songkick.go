package songkick

import (
	"net/http"
	"log"
	"encoding/json"
	"time"
	"fmt"
)

type SongkickClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func (c *SongkickClient) GetUpcoming(locationId, months int) []Event {

	endDate := time.Now().AddDate(0, 0, 7)
	log.Printf(`end date is %s`, endDate.Format(`2006-01-02`))

	events := make([]Event, 0)
	page := 1
	for true {

		log.Printf(`loading page %d`, page)

		request, _ := http.NewRequest(`GET`, `http://api.songkick.com/api/3.0/metro_areas/`+fmt.Sprintf(`%d`, locationId)+`/calendar.json?apikey=`+c.APIKey+`&page=`+fmt.Sprintf(`%d`, page), nil)
		request.Header.Set(`User-Agent`, `BerlinBot/0.1 by warmans`)

		res, err := c.HTTPClient.Do(request)
		if err != nil {
			log.Printf("Failed to rerieve calender because: %s", err.Error())
			return events
		}

		result := &SongkickResult{}
		decodeErr := json.NewDecoder(res.Body).Decode(result)
		res.Body.Close()

		if decodeErr != nil {
			log.Printf("Failed to decode songkick response: %s", decodeErr.Error())
			return events
		}

		for _, event := range result.Resultspage.Results.Event {
			if eventDate, err := time.Parse("2006-01-02", event.Start.Date); err == nil {
				if eventDate.Unix() < endDate.Unix() {
					events = append(events, event)
				} else {
					return events //we encoutered a date outside of the required window so we can return
				}
			} else {
				log.Printf(`Failed to parse event date %s`, err.Error())
			}
		}

		//move to next page
		page++;

	}

	return events
}


type SongkickResult struct {
	Resultspage struct {
		Status string `json:"status"`
		Results struct {
			Event []Event `json:"event"`
		} `json:"results"`
		Perpage int `json:"perPage"`
		Page int `json:"page"`
		Totalentries int `json:"totalEntries"`
	} `json:"resultsPage"`
}

type Event struct {
	Type string `json:"type"`
	Displayname string `json:"displayName"`
	Status string `json:"status"`
	Popularity float64 `json:"popularity"`
	Start struct {
		Time string `json:"time"`
		Datetime string `json:"datetime"`
		Date string `json:"date"`
	} `json:"start"`
	Agerestriction interface{} `json:"ageRestriction"`
	Location struct {
		City string `json:"city"`
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	URI string `json:"uri"`
	Performance []struct {
		Billing string `json:"billing"`
		Billingindex int `json:"billingIndex"`
		Displayname string `json:"displayName"`
		Artist struct {
			Identifier []struct {
				Href string `json:"href"`
				Mbid string `json:"mbid"`
			} `json:"identifier"`
			Displayname string `json:"displayName"`
			URI string `json:"uri"`
			ID int `json:"id"`
		} `json:"artist"`
		ID int `json:"id"`
	} `json:"performance"`
	Venue struct {
		Displayname string `json:"displayName"`
		Metroarea struct {
			Displayname string `json:"displayName"`
			URI string `json:"uri"`
			ID int `json:"id"`
			Country struct {
				Displayname string `json:"displayName"`
			} `json:"country"`
		} `json:"metroArea"`
		Lat float64 `json:"lat"`
		URI string `json:"uri"`
		Lng float64 `json:"lng"`
		ID int `json:"id"`
	} `json:"venue"`
	ID int `json:"id"`

}