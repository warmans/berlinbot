Concerts this Month
======================

The following concerts are happening this month.

| Date | Event | Venue | Artist(s) | Info |
|------|-------|-------|-----------|------|
{{range .}}|{{.Start.Date}} | {{.Displayname}} | {{.Venue.Displayname}} | {{range .Performance}}[{{.Artist.Displayname}}](spotify:search:{{.Artist.Displayname}}),{{end}} | [Songkick]({{.URI}})|{{end}}
