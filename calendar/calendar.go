// Package calendar contains the high level types and functions to handle
// iCalendar files.
package calendar

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"home.leo-peltier.fr/goof/calendar/ics"
)

// A Calendar is an ICS VCALENDAR.
type Calendar struct {
	Events []Event
}

// An Event is an ICS VEVENT.
// See RFC 5545 3.6.1.
type Event struct {
	Summary     string
	Description string
	Location    string
	UID         string
	Start       time.Time
	End         time.Time
	Transparent bool

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
		state   int = STATE_INIT
		event   Event
		version string
	)

loop:
	for scanner.Scan() {
		line := ics.NewLine(scanner.Text())

		switch state {
		case STATE_INIT:
			// BEGIN:VCALENDAR being the first line is mandatory per RFC.
			if line.Begins(ics.CalendarComponent) {
				state = STATE_PARSE_CALENDAR
			} else {
				errs = append(errs, errors.New("Expected BEGIN:VCALENDAR"))
				state = STATE_END
			}
		case STATE_PARSE_CALENDAR:
			// Chomp until we read an event or exit the calendar.
			if line.Begins(ics.EventComponent) {
				event = Event{}
				state = STATE_PARSE_EVENT
			} else if line.Property == ics.VersionProperty {
				version = line.Value
			} else if line.Ends(ics.CalendarComponent) {
				state = STATE_END
			}
		case STATE_PARSE_EVENT:
			if line.Begins(ics.AlarmComponent) {
				state = STATE_PARSE_ALARM
			} else if line.Ends(ics.EventComponent) {
				event = event.InferMissingValues()
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
			if line.Ends(ics.AlarmComponent) {
				state = STATE_PARSE_EVENT
			}
		case STATE_END:
			break loop
		default:
			panic("unreachable")
		}
	}

	if version != "2.0" {
		errs = append(errs, errors.New("Version is not '2.0', may not be an iCalendar."))
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

// FromFile read an iCalendar from a path.
func FromFile(path string) (calendar Calendar, errs []error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		errs = append(errs, err)
		return
	}

	return FromReader(file)
}

// UpdateFromIcsProperty sets an event property from an ICS line.
func (e *Event) UpdateFromIcsLine(line ics.Line) (err error) {
	switch line.Property {
	case ics.DateTimeStartProperty:
		e.Start, err = line.ParseAsTime()
	case ics.DateTimeEndProperty:
		e.End, err = line.ParseAsTime()
	case ics.SummaryProperty:
		e.Summary = line.Value
	case ics.DescriptionProperty:
		e.Description = line.Value
	case ics.LocationProperty:
		e.Location = line.Value
	case ics.CreatedProperty:
		e.Created, err = line.ParseAsTime()
	case ics.LastModifiedProperty:
		e.LastModified, err = line.ParseAsTime()
	case ics.UIDProperty:
		e.UID = line.Value
	case ics.TransparentProperty:
		if line.Value == ics.TransparentTransparency {
			e.Transparent = true
		} else if line.Value == ics.OpaqueTransparency {
			e.Transparent = false
		} else {
			err = errors.New("Invalid transparency: " + line.Value)
		}
	}

	return err
}

func (orig Event) InferMissingValues() Event {
	e := orig

	// Create 0-duration events when end/start is missing.
	if !orig.Start.IsZero() && orig.End.IsZero() {
		e.End = orig.Start
	} else if orig.Start.IsZero() && !orig.End.IsZero() {
		e.Start = orig.End
	}

	if len(e.UID) == 0 {
		hash := sha1.New()
		io.WriteString(hash, e.Start.String())
		io.WriteString(hash, e.Summary)
		e.UID = fmt.Sprintf("%x@goof", hash.Sum(nil))
	}

	return e
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

type CalendarFilter struct {
	RangeUpper time.Time
	RangeLower time.Time
}

func (original Calendar) Filter(filter CalendarFilter) (filtered Calendar) {
	for _, event := range original.Events {
		startsInRange := event.Start.After(filter.RangeLower) && event.Start.Before(filter.RangeUpper)
		endsInRange := event.End.After(filter.RangeLower) && event.End.Before(filter.RangeUpper)
		if startsInRange || endsInRange {
			filtered.Events = append(filtered.Events, event)
		}
	}
	return
}
