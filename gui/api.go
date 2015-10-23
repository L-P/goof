package gui

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji/web"

	"home.leo-peltier.fr/goof/calendar"
)

func getCalendar(c web.C, r *http.Request) (data interface{}, err error) {
	calendars := c.Env["calendars"].(map[string]calendar.Calendar)
	if _, prs := calendars[c.URLParams["calendar"]]; !prs {
		err = errors.New("Calendar not found.")
		return
	}

	var (
		fullCalendar     calendar.Calendar = calendars[c.URLParams["calendar"]]
		filteredCalendar calendar.Calendar
		filter           calendar.CalendarFilter
	)

	r.ParseForm()
	if r.Form.Get("range") != "" {
		filter.RangeLower, filter.RangeUpper, err = parseRange(r.Form.Get("range"))
		if err != nil {
			return
		}
	}

	filteredCalendar, err = fullCalendar.Filter(filter)
	if err != nil {
		return
	}

	data = struct{ Calendar calendar.Calendar }{
		Calendar: filteredCalendar,
	}

	return data, err
}

// parseRange takes a string (eg. 2006-01-02,2006-02-02) and return the parsed times.
func parseRange(str string) (lower, upper time.Time, err error) {
	splits := strings.SplitN(str, ",", 2)
	if len(splits) != 2 {
		err = errors.New("Bad value for 'range' parameter.")
		return
	}

	if splits[0] != "" {
		lower, err = time.Parse("2006-01-02", splits[0])
		if err != nil {
			return
		}
	}

	if splits[1] != "" {
		upper, err = time.Parse("2006-01-02", splits[1])
		if err != nil {
			return
		}
	}

	if !upper.IsZero() && !lower.IsZero() {
		if upper.Before(lower) {
			err = errors.New("upper < lower")
			return
		}
	}

	return
}
