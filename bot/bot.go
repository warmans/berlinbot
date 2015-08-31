package bot

import (
	"log"
	"bytes"
	"github.com/warmans/berlinbot/songkick"
	"text/template"
	"github.com/warmans/berlinbot/spotify"
	"fmt"
	"strings"
)

type Bot struct {
	Debug    bool
	Reddit   *RedditClient
	Songkick *songkick.SongkickClient
	Spotify  *spotify.Spotify
}

func (b *Bot) SubmitEvents(sub, title string) {

	helpers := template.FuncMap{
		"join": func(vals []string, seperator string) string {
			out := strings.TrimSpace(strings.Join(vals, seperator))
			if out == "" {
				return "NA"
			}
			return out
		},
		"link_map": func(linkMap map[string]string) string {
			var out string
			for name, uri := range linkMap {
				out += fmt.Sprintf("[%s](%s) ", name, uri)
			}
			return out
		},
	}

	tpl, err := template.New("post.tpl.md").Funcs(helpers).ParseFiles(`post.tpl.md`)

	if err != nil {
		log.Printf(`Failed to render template: %s`, err.Error())
		return
	}

	viewModel := &EventsViewModel{spotify: b.Spotify}
	viewModel.ImportEvents(b.Songkick.GetUpcoming(28443, 1))

	out := bytes.NewBufferString(``)
	if err := tpl.Execute(out, viewModel); err != nil {
		log.Printf(`failed to execute template: %s`, err.Error())
	}

	if b.Debug {
		fmt.Println(title)
		fmt.Println(out.String())
	} else {
		b.Reddit.Submit(SUBMIT_TYPE_SELF, title, out.String(), sub)
	}
}
