package main

import (
	"encoding/json"
	"fmt"
	"os"

	"home.leo-peltier.fr/goof/calendar"
)

func main() {
	if len(os.Args) != 2 {
		panic("Need one .ics as sole argument.")
	}

	ics, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	cal, errs := calendar.FromReader(ics)
	if len(errs) > 0 {
		fmt.Println("Errors occured when parsing the iCalendar:")
		fmt.Println(errs)
	}

	out, _ := json.Marshal(cal.Events)
	fmt.Printf("%s\n", out)
}
