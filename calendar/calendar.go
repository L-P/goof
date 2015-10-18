// calendar contains the high level strucs functions to handle iCalendar files.
package calendar

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"time"

	"./ics"
)

// An Event is an ICS VCALENDAR.
type Calendar struct {
	Events []Event
}

// An Event is an ICS VEVENT.
type Event struct {
	Summary     string
	Description string
	UID         string
	Start       time.Time
	Duration    time.Duration
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
func FromReader(reader io.Reader) (calendar Calendar, err error) {
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
				err = errors.New("Expected BEGIN:VCALENDAR")
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
				event.UpdateFromIcsLine(line)
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
		err = errors.New("Parsing failed to end correctly.")
	} else if scanner.Err() != nil {
		err = scanner.Err()
	}

	sort.Sort(ByStart(calendar.Events))

	return
}

// UpdateFromIcsProperty sets an event property from an ICS line.
func (e *Event) UpdateFromIcsLine(line ics.Line) {
	switch line.Property {
	case "DTSTART":
		e.Start = parseTime(line)
	case "SUMMARY":
		e.Summary = line.Value
	case "DESCRIPTION":
		e.Description = line.Value
	case "UID":
		e.UID = line.Value
	}
}

func parseTime(line ics.Line) time.Time {
	valueType, _ := line.Parameters["VALUE"]

	var (
		parsed time.Time
		err    error
	)

	if valueType == "DATE" {
		parsed, err = time.Parse("20060102", line.Value)
	} else if line.Value[len(line.Value)-1] == 'Z' {
		parsed, err = time.Parse("20060102T150405Z", line.Value)
	} else {
		valueTz, hasTz := line.Parameters["TZID"]
		var loc *time.Location = time.UTC
		if hasTz {
			loc, _ = time.LoadLocation(valueTz)
		}
		parsed, err = time.ParseInLocation("20060102T150405", line.Value, loc)
	}

	if err != nil {
		panic(err)
	}

	return parsed.In(time.UTC)
}

type ByStart []Event

func (t ByStart) Len() int {
	return len(t)
}

func (t ByStart) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByStart) Less(i, j int) bool {
	return t[i].Start.Unix() < t[j].Start.Unix()
}
