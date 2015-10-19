package gui

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

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

func sendJSONResponse(c web.C, w http.ResponseWriter, data interface{}) {
	var response JSONResponse
	response.Data = data
	response.Meta.Errors = c.Env["errors"].([]error)

	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Meta.Errors)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func Serve() {
	goji.Use(loadCalendars)

	goji.Get("/", handler)
	goji.Get("/calendar/:calendar", calendarHandler)

	goji.DefaultMux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	goji.Serve()
}

func calendarHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	data, err := getCalendar(c, r)
	errs := make([]error, 0)
	if err != nil {
		errs = append(errs, err)
	}
	c.Env["errors"] = errs
	sendJSONResponse(c, w, data)
}

func parseRange(str string) (lower, upper time.Time, err error) {
	splits := strings.SplitN(str, ",", 2)
	if len(splits) != 2 {
		err = errors.New("Bad value for 'range' parameter.")
		return
	}

	lower, err = time.Parse("2006-01-02", splits[0])
	if err != nil {
		return
	}

	upper, err = time.Parse("2006-01-02", splits[1])
	if err != nil {
		return
	}

	if upper.Before(lower) {
		err = errors.New("upper < lower")
		return
	}

	return
}

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

func getCalendar(c web.C, r *http.Request) (data struct{ Calendar calendar.Calendar }, err error) {
	calendars := c.Env["calendars"].(map[string]calendar.Calendar)
	if _, prs := calendars[c.URLParams["calendar"]]; !prs {
		err = errors.New("Calendar not found.")
		return
	}

	fullCalendar := calendars[c.URLParams["calendar"]]

	r.ParseForm()
	if r.Form.Get("range") == "" {
		err = errors.New("Range parameter is mandatory.")
		return
	}

	lower, upper, err := parseRange(r.Form.Get("range"))
	if err != nil {
		return
	}

	if upper.Sub(lower).Seconds() > 2*3600*24*31 {
		err = errors.New("Range > 2 month.")
		return
	}

	data.Calendar = fullCalendar.Filter(
		calendar.CalendarFilter{
			RangeUpper: upper,
			RangeLower: lower,
		},
	)

	return
}
