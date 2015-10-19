package gui

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"home.leo-peltier.fr/goof/calendar"
)

type JSONResponse struct {
	Data interface{}
	Meta struct {
		Errors []error
	}
}

func handler(c web.C, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("templates/*.html"))
	t.ExecuteTemplate(w, "index", nil)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, errs []error) {
	var response JSONResponse
	response.Data = data
	response.Meta.Errors = errs

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getCalendar(c web.C, w http.ResponseWriter, r *http.Request) {
	data := struct {
		Calendars []calendar.Calendar
	}{
		make([]calendar.Calendar, 0),
	}
	errs := make([]error, 0)

	// TODO: calendars path should exist in a configuration file.
	cal, errs := calendar.FromFile(c.URLParams["name"])
	data.Calendars = append(data.Calendars, cal)
	sendJSONResponse(w, data, errs)
}

func Serve() {
	goji.Get("/", handler)
	goji.Get("/calendar/:name", getCalendar)

	goji.DefaultMux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	goji.Serve()
}
