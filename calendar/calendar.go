package calendar

import (
	"errors"
	"os"
	"time"
)

type Calendar struct {
	Events []Event
}

type Event struct {
	Summary     string
	Description string
	UID         string
	Time        time.Time
	Duration    time.Duration
}

const (
	STATE_INIT int = iota
	STATE_END
	STATE_PARSE_CALENDAR
	STATE_PARSE_EVENT
)

func FromFile(file *os.File) (calendar Calendar, err error) {
	scanner := NewScanner(file)
	var (
		state int = STATE_INIT
		event Event
	)

loop:
	for scanner.Scan() {
		key, value := scanner.KeyValue()

		switch state {
		case STATE_INIT:
			// Make BEGIN:VCALENDAR the mandatory first line of an .ics.
			if key == "BEGIN" && value == "VCALENDAR" {
				state = STATE_PARSE_CALENDAR
			} else {
				err = errors.New("Expected BEGIN:VCALENDAR")
				state = STATE_END
			}
		case STATE_PARSE_CALENDAR:
			// Chomp until we read an event or exit the calendar.
			if key == "BEGIN" && value == "VEVENT" {
				event = Event{}
				state = STATE_PARSE_EVENT
			} else if key == "END" && value == "VCALENDAR" {
				state = STATE_END
			}
		case STATE_PARSE_EVENT:
			if key == "END" && value == "VEVENT" {
				calendar.Events = append(calendar.Events, event)
				state = STATE_PARSE_CALENDAR
			} else {
				event.UpdateFromIcsProperty(key, value)
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

func (e *Event) UpdateFromIcsProperty(name string, value string) {
	// TODO
}
