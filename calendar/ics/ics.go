package ics

import (
	"errors"
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
	case 0:
		err = errors.New("No ':' found in line.")
		return
	case 1:
		line.Property, line.Parameters = parseProperty(splits[0])
	case 2:
		line.Property, line.Parameters = parseProperty(splits[0])
		line.Value = splits[1]
	default:
		panic("unreachable")
	}

	line.Property = splits[0]
	return
}

func parseProperty(str string) (property string, parameters map[string]string) {
	return
}
