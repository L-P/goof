// Package ics contains the utilities for parsing iCalendar files.
package ics

import (
	"strings"
	"time"
)

// Line contains the parsed data of a single line from an iCalendar file.
// Every property and parameter names will be uppercased.
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
	line.Property, line.Parameters = parseProperty(splits[0])

	if len(splits) > 1 {
		line.Value = splits[1]
	}

	return line
}

// parseProperty parse the property name, its parameters and their
// values from first part of an iCalendar line (before the first ':').
func parseProperty(str string) (property string, parameters map[string]string) {
	splits := strings.SplitN(str, ";", 2)

	property = strings.ToUpper(splits[0])
	if len(splits) > 1 {
		parameters = parsePropertyParameters(splits[1])
	}

	return property, parameters
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

// ParseAsTime parses the value of a line as a date or timestamp.
// Time will be given in its original timezone if given in the TZID
// parameter or UTC by default.
// See RFC 5545 3.3.4. and 3.3.5.
func (line Line) ParseAsTime() (parsed time.Time, err error) {
	valueType, _ := line.Parameters["VALUE"]

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

	return
}
