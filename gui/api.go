package gui

import (
	"errors"
	"net/http"
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
	filter.RangeLower, err = time.Parse("2006-01-02", r.Form.Get("start"))
	if err != nil {
		return
	}

	filter.RangeUpper, err = time.Parse("2006-01-02", r.Form.Get("end"))
	if err != nil {
		return
	}

	if !filter.RangeUpper.IsZero() && !filter.RangeLower.IsZero() {
		if filter.RangeUpper.Before(filter.RangeLower) {
			err = errors.New("upper < lower")
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
