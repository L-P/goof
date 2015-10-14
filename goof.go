package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"./calendar"
)

func handler(c web.C, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("templates/*.html"))
	days := []calendar.Day{
		{
			Date: time.Now(),
			Events: []calendar.Event{
				{
					Time:        time.Now(),
					Title:       "test event",
					Description: "longer desc",
				},
			},
		},
	}

	t.ExecuteTemplate(w, "index", struct {
		days []calendar.Day
	}{
		days: days,
	})

	fmt.Println(days)
}

func main() {
	goji.Get("/", handler)
	goji.DefaultMux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	goji.Serve()
}
