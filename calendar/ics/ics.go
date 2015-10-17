package ics

import (
	"strings"
)

type Line struct {
	Property   string
	Parameters map[string]string
	Value      string

	original string
}

func (l Line) String() string {
	return l.original
}

func NewLine(str string) (line Line, err error) {
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
