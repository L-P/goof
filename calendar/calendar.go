package calendar

import (
	"time"
)

type Event struct {
	Time        time.Time
	Duration    time.Duration
	Title       string
	Description string
}

type Day struct {
	Date   time.Time
	Events []Event
}
