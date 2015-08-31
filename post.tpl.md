| Popularity | Date | Genre(s) | Event | Venue | Artist(s) | Info |
|------------|------|----------|-------|-------|-----------|------|
{{range .Events}}| {{printf "%.2f" .Popularity}} | {{.Date}} | {{join .Genres ", "}} | {{.Name}} | {{.Venue}} | {{join .Artists ", "}} | {{link_map .Links}} |
{{end}}
