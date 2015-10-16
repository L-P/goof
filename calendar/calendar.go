// calendar contains the high level strucs functions to handle iCalendar files.
package calendar

import (
	"bufio"
	"errors"
	"io"
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
		event Event
	)

loop:
	for scanner.Scan() {
		line := ics.NewLine(scanner.Text())

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

	return
}

// UpdateFromIcsProperty sets an event property from an ICS line.
func (e *Event) UpdateFromIcsLine(line ics.Line) {
	// TODO
}
