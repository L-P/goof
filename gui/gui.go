package gui

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"home.leo-peltier.fr/goof/calendar"
)

type ApiEndpoint func(c web.C, r *http.Request) (data interface{}, err error)

// apiHandler wraps an ApiEndpoint to add error and content-type handling.
func apiHandler(apiFunc ApiEndpoint) web.HandlerType {
	wrap := func(c web.C, w http.ResponseWriter, r *http.Request) {
		data, err := apiFunc(c, r)
		errs := make([]error, 0)
		if err != nil {
			errs = append(errs, err)
		}
		c.Env["errors"] = errs
		sendJSONResponse(c, w, data)
	}

	return web.HandlerFunc(wrap)
}

type JSONResponse struct {
	Data interface{}
	Meta struct {
		Errors []string
	}
}

func sendJSONResponse(c web.C, w http.ResponseWriter, data interface{}) {
	var response JSONResponse
	response.Data = data
	for _, err := range c.Env["errors"].([]error) {
		response.Meta.Errors = append(response.Meta.Errors, err.Error())
	}

	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func Serve() {
	goji.Use(loadCalendars)

	goji.Get("/", rootHandler)
	goji.Get("/calendar/:calendar", apiHandler(getCalendar))

	goji.DefaultMux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	goji.Serve()
}

func rootHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	t := template.New("")
	// Solve mustache/go delimiter conflict by using <% %> in views.
	t.Delims("<%", "%>").ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "index", nil)
}

// loadCalendars is a Goji middleware that injects calendars in the request
// context.
func loadCalendars(c *web.C, h http.Handler) http.Handler {
	errs := make([]error, 0)
	calendars := make(map[string]calendar.Calendar, 0)
	calendars["calendar.ics"], errs = calendar.FromFile("calendar.ics")

	if len(errs) > 0 {
		panic("Unable to load calendars.")
	}

	wrap := func(w http.ResponseWriter, r *http.Request) {
		if c.Env == nil {
			c.Env = make(map[interface{}]interface{})
		}

		c.Env["calendars"] = calendars
		c.Env["errors"] = errs

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(wrap)
}
