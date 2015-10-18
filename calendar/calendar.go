// Package calendar contains the high level types and functions to handle
// iCalendar files.
package calendar

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"time"

	"home.leo-peltier.fr/goof/calendar/ics"
)

// An Event is an ICS VCALENDAR.
type Calendar struct {
	Events []Event
}

// An Event is an ICS VEVENT.
type Event struct {
	Summary     string
	Description string
	Location    string
	UID         string
	Start       time.Time
	Duration    time.Duration

	Created      time.Time
	LastModified time.Time
}

// States for the parser in FromReader which is implemented as a FSM.
const (
	STATE_INIT int = iota
	STATE_END
	STATE_PARSE_CALENDAR
	STATE_PARSE_EVENT
	STATE_PARSE_ALARM
)

// FromReader reads ICS data and create a Calendar filled with Events.
func FromReader(reader io.Reader) (calendar Calendar, errs []error) {
	scanner := bufio.NewScanner(reader)
	var (
		state int = STATE_INIT
		event Event
	)

loop:
	for scanner.Scan() {
		line := ics.NewLine(scanner.Text())

		switch state {
		case STATE_INIT:
			// BEGIN:VCALENDAR being the first line is mandatory per RFC.
			if line.BeginsCalendar() {
				state = STATE_PARSE_CALENDAR
			} else {
				errs = append(errs, errors.New("Expected BEGIN:VCALENDAR"))
				state = STATE_END
			}
		case STATE_PARSE_CALENDAR:
			// Chomp until we read an event or exit the calendar.
			if line.BeginsEvent() {
				event = Event{}
				state = STATE_PARSE_EVENT
			} else if line.EndsCalendar() {
				state = STATE_END
			}
		case STATE_PARSE_EVENT:
			if line.BeginsAlarm() {
				state = STATE_PARSE_ALARM
			} else if line.EndsEvent() {
				calendar.Events = append(calendar.Events, event)
				state = STATE_PARSE_CALENDAR
			} else {
				err := event.UpdateFromIcsLine(line)
				if err != nil {
					errs = append(errs, err)
				}
			}
		case STATE_PARSE_ALARM:
			// Chomp until we get back to the event.
			if line.EndsAlarm() {
				state = STATE_PARSE_EVENT
			}
		case STATE_END:
			break loop
		default:
			panic("unreachable")
		}
	}

	if state != STATE_END {
		errs = append(errs, errors.New("Parsing failed to end correctly."))
	}
	if scanner.Err() != nil {
		errs = append(errs, scanner.Err())
	}

	sort.Sort(byStart(calendar.Events))

	return calendar, errs
}

// UpdateFromIcsProperty sets an event property from an ICS line.
func (e *Event) UpdateFromIcsLine(line ics.Line) (err error) {
	switch line.Property {
	case "DTSTART":
		var parsed time.Time
		parsed, err = line.ParseAsTime()
		e.Start = parsed
	case "SUMMARY":
		e.Summary = line.Value
	case "DESCRIPTION":
		e.Description = line.Value
	case "LOCATION":
		e.Location = line.Value
	case "CREATED":
		var parsed time.Time
		parsed, err = line.ParseAsTime()
		e.Created = parsed
	case "LAST-MODIFIED":
		var parsed time.Time
		parsed, err = line.ParseAsTime()
		e.LastModified = parsed
	case "UID":
		e.UID = line.Value
	}

	return err
}

type byStart []Event

func (t byStart) Len() int {
	return len(t)
}

func (t byStart) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byStart) Less(i, j int) bool {
	return t[i].Start.Unix() < t[j].Start.Unix()
}
