package bot

import (
	"github.com/warmans/berlinbot/songkick"
	"github.com/warmans/berlinbot/spotify"
	"math"
)

//EventsViewModel allows for a simple template
type EventsViewModel struct {
	spotify *spotify.Spotify
	Events  []*Event
}

func (vm *EventsViewModel) ImportEvents(SKEvents []songkick.Event) {

	vm.Events = make([]*Event, len(SKEvents))
	for k, ev := range SKEvents {

		//basic reliable info
		vm.Events[k] = &Event{
			Date: ev.Start.Date,
			Name: ev.Displayname,
			Venue: ev.Venue.Displayname,
			Popularity: 0.0,
			Artists: make([]string, len(ev.Performance)),
			Genres: make([]string, 0),
			Links: map[string]string{"Songkick": ev.URI},
		}

		//varied artist info
		popularityDivisor := 0.0
		for ak, songkickArt := range ev.Performance {
			vm.Events[k].Artists[ak] = songkickArt.Displayname
			spotifyArt := vm.spotify.GetArtistByName(songkickArt.Displayname)
			if spotifyArt != nil {

				//append genres where available
				if len(spotifyArt.Genres) > 0 {
					vm.Events[k].Genres = append(vm.Events[k].Genres, spotifyArt.Genres[0])
				}

				//increase songkick popularity by spotify popularity
				vm.Events[k].Popularity += float64(spotifyArt.Popularity)
				popularityDivisor++

				//add a link to the artist's page on spotify
				vm.Events[k].Links[songkickArt.Displayname] = spotifyArt.ExternalUrls.Spotify
			}
		}

		//reduce implact of large artist lists and boost by sonkick popularity
		vm.Events[k].Popularity = (vm.Events[k].Popularity / math.Min(1.0, popularityDivisor)) * math.Min(ev.Popularity, 1.0)
		if math.IsNaN(vm.Events[k].Popularity) {
			vm.Events[k].Popularity = 0.0
		}
	}
}

type Event struct {
	Date       string
	Name       string
	Venue      string
	Popularity float64
	Artists    []string
	Genres     []string
	Links      map[string]string
}
