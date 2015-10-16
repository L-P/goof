package calendar

import (
	"fmt"
	"io"
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

	Created      time.Time
	LastModified time.Time
}

func FromFile(file *os.File) (calendar Calendar, err error) {
	reader := NewScanner(file)
	var key, value string

	for err == nil {
		key, value, err = reader.Next()
		fmt.Printf("k: %s, v:%s\n", key, value)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
