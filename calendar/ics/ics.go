package ics

import (
	"strings"
)

// Line contains the parsed data of a single line from an iCalendar file.
type Line struct {
	Property   string
	Parameters map[string]string
	Value      string

	original string
}

// BeginsCalendar returns true if the line is starting a new VCALENDAR.
func (l Line) BeginsCalendar() bool {
	return l.Property == "BEGIN" && l.Value == "VCALENDAR"
}

// EndsCalendar returns true if the line is ending the current VCALENDAR.
func (l Line) EndsCalendar() bool {
	return l.Property == "END" && l.Value == "VCALENDAR"
}

// BeginsEvent returns true if the line is starting a new VEVENT.
func (l Line) BeginsEvent() bool {
	return l.Property == "BEGIN" && l.Value == "VEVENT"
}

// EndsEvent returns true if the line is ending the current VEVENT.
func (l Line) EndsEvent() bool {
	return l.Property == "END" && l.Value == "VEVENT"
}

// BeginsAlarm returns true if the line is starting a new VALARM.
func (l Line) BeginsAlarm() bool {
	return l.Property == "BEGIN" && l.Value == "VALARM"
}

// EndsAlarm returns true if the line is ending the current VALARM.
func (l Line) EndsAlarm() bool {
	return l.Property == "END" && l.Value == "VALARM"
}

func (l Line) String() string {
	return l.original
}

func NewLine(str string) (line Line) {
	line.original = str

	splits := strings.SplitN(str, ":", 2)
	switch len(splits) {
	case 1:
		line.Property, line.Parameters = parseProperty(splits[0])
	case 2:
		line.Property, line.Parameters = parseProperty(splits[0])
		line.Value = splits[1]
	default:
		panic("unreachable")
	}

	return
}

// parseProperty parse the property name, its parameters and their
// values from first part of an iCalendar line (before the first ':').
func parseProperty(str string) (property string, parameters map[string]string) {
	splits := strings.SplitN(str, ";", 2)

	switch len(splits) {
	case 1:
		property = splits[0]
	case 2:
		property = splits[0]
		parameters = parsePropertyParameters(splits[1])
	default:
		panic("unreachable")
	}

	property = strings.ToUpper(property)

	return
}

// parsePropertyParameters return a property parameters and their values
// from a string.
// From the line "PROP;PARAM=VALUE;OTHER=VALUE", only pass the part
// after "PROP;"
func parsePropertyParameters(str string) (parameters map[string]string) {
	parameters = make(map[string]string)
	tuples := strings.Split(str, ";")
	for _, tuple := range tuples {
		splits := strings.SplitN(tuple, "=", 2)
		var value string
		if len(splits) > 1 {
			value = splits[1]
		}
		parameters[strings.ToUpper(splits[0])] = value
	}

	return
}
