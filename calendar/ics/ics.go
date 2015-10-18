// Package ics contains the utilities for parsing iCalendar files.
package ics

import (
	"strings"
	"time"
)

// Parameter is a parameter name. See RFC 5545 8.3.3.
type Parameter string

const (
	ValueTypeParameter  Parameter = "VALUE"
	TimezoneIDParameter           = "TZID"
)

// Values for ValueTypeParameter.
const (
	DateType     = "DATE"
	DateTimeType = "DATE-TIME"
)

// Property is a property name. See RFC 5545 8.3.2.
type Property string

const (
	BeginProperty         Property = "BEGIN"
	EndProperty                    = "END"
	VersionProperty                = "VERSION"
	DateTimeStartProperty          = "DTSTART"
	DateTimeEndProperty            = "DTEND"
	SummaryProperty                = "SUMMARY"
	DescriptionProperty            = "DESCRIPTION"
	LocationProperty               = "LOCATION"
	CreatedProperty                = "CREATED"
	LastModifiedProperty           = "LAST-MODIFIED"
	UIDProperty                    = "UID"
	TransparentProperty            = "TRANSP"
)

// Values for TransparentProperty.
const (
	TransparentTransparency = "TRANSPARENT"
	OpaqueTransparency      = "OPAQUE"
)

// Component is a component name. See RFC 5545 8.3.1.
type Component string

const (
	CalendarComponent Component = "VCALENDAR"
	EventComponent              = "VEVENT"
	AlarmComponent              = "VALARM"
	TimezoneComponent           = "VTIMEZONE"
	JournalComponent            = "VJOURNAL"
	TodoComponent               = "VTODO"
	FreeBusyComponent           = "VFREEBUZY"
)

// Line contains the parsed data of a single line from an iCalendar file.
// Every property and parameter names will be uppercased.
type Line struct {
	Property   Property
	Parameters map[Parameter]string
	Value      string

	original string
}

// Begins returns true if the line is starting a new component.
func (l Line) Begins(component Component) bool {
	return l.Property == BeginProperty && l.Value == string(component)
}

// Ends returns true if the line is ending the current component.
func (l Line) Ends(component Component) bool {
	return l.Property == EndProperty && l.Value == string(component)
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
func parseProperty(str string) (property Property, parameters map[Parameter]string) {
	splits := strings.SplitN(str, ";", 2)

	property = Property(strings.ToUpper(splits[0]))
	if len(splits) > 1 {
		parameters = parsePropertyParameters(splits[1])
	}

	return property, parameters
}

// parsePropertyParameters return a property parameters and their values
// from a string.
// From the line "PROP;PARAM=VALUE;OTHER=VALUE", only pass the part
// after "PROP;"
func parsePropertyParameters(str string) (parameters map[Parameter]string) {
	parameters = make(map[Parameter]string)
	tuples := strings.Split(str, ";")
	for _, tuple := range tuples {
		splits := strings.SplitN(tuple, "=", 2)
		var value string
		if len(splits) > 1 {
			value = splits[1]
		}

		name := Parameter(strings.ToUpper(splits[0]))
		parameters[name] = value
	}

	return
}

// ParseAsTime parses the value of a line as a date or timestamp.
// Time will be given in its original timezone if given in the TZID
// parameter or UTC by default.
// See RFC 5545 3.3.4. and 3.3.5.
func (line Line) ParseAsTime() (parsed time.Time, err error) {
	valueType, _ := line.Parameters[ValueTypeParameter]

	if valueType == DateType {
		parsed, err = time.Parse("20060102", line.Value)
	} else if line.Value[len(line.Value)-1] == 'Z' {
		parsed, err = time.Parse("20060102T150405Z", line.Value)
	} else {
		valueTz, hasTz := line.Parameters[TimezoneIDParameter]
		var loc *time.Location = time.UTC
		if hasTz {
			loc, _ = time.LoadLocation(valueTz)
		}
		parsed, err = time.ParseInLocation("20060102T150405", line.Value, loc)
	}

	return parsed, err
}
