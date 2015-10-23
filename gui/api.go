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

	data = struct {
		Calendar calendar.Calendar
	}{
		Calendar: fullCalendar.Filter(
			calendar.CalendarFilter{
				RangeUpper: upper,
				RangeLower: lower,
			},
		),
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
