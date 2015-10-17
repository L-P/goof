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
	Time        time.Time
	Duration    time.Duration
}

// States for the parser in FromReader which is implemented as a FSM.
const (
	STATE_INIT int = iota
	STATE_END
	STATE_PARSE_CALENDAR
	STATE_PARSE_EVENT
)

// FromReader reads ICS data and create a Calendar filled with Events.
func FromReader(reader io.Reader) (calendar Calendar, err error) {
	scanner := bufio.NewScanner(reader)
	var (
		state int = STATE_INIT
		line  ics.Line
		event Event
	)

loop:
	for scanner.Scan() {
		line, err = ics.NewLine(scanner.Text())

		if err != nil {
			return
		}

		switch state {
		case STATE_INIT:
			// Make BEGIN:VCALENDAR the mandatory first line of an .ics.
			if line.String() == "BEGIN:VCALENDAR" {
				state = STATE_PARSE_CALENDAR
			} else {
				err = errors.New("Expected BEGIN:VCALENDAR")
				state = STATE_END
			}
		case STATE_PARSE_CALENDAR:
			// Chomp until we read an event or exit the calendar.
			if line.String() == "BEGIN:VEVENT" {
				event = Event{}
				state = STATE_PARSE_EVENT
			} else if line.String() == "END:VCALENDAR" {
				state = STATE_END
			}
		case STATE_PARSE_EVENT:
			if line.String() == "END:VEVENT" {
				calendar.Events = append(calendar.Events, event)
				state = STATE_PARSE_CALENDAR
			} else {
				event.UpdateFromIcsLine(line)
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

	sort.Sort(ByTime(calendar.Events))

	return
}

// UpdateFromIcsProperty sets an event property from an ICS line.
func (e *Event) UpdateFromIcsLine(line ics.Line) {
	switch line.Property {
	case "DTSTART":
		e.Time = parseTime(line)
	case "SUMMARY":
		e.Summary = line.Value
	case "DESCRIPTION":
		e.Description = line.Value
	}
}

type ByTime []Event

func (t ByTime) Len() int {
	return len(t)
}

func (t ByTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByTime) Less(i, j int) bool {
	return t[i].Time.Unix() < t[j].Time.Unix()
}

func parseTime(line ics.Line) time.Time {
	valueType, prs := line.Parameters["VALUE"]
	var parsed time.Time
	var err error

	if prs && valueType == "DATE" {
		parsed, err = time.Parse("20060102", line.Value)
	} else if line.Value[len(line.Value)-1] == 'Z' {
		parsed, err = time.Parse("20060102T150405Z", line.Value)
	} else {
		parsed, err = time.Parse("20060102T150405", line.Value)
	}

	if err != nil {
		panic(err)
	}
	return parsed
}
