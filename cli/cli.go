package main

import (
	"fmt"
	"os"

	"../calendar"
)

func main() {
	if len(os.Args) != 2 {
		panic("Need one .ics as sole argument.")
	}

	ics, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	cal, err := calendar.FromFile(ics)
	if err != nil {
		panic(err)
	}

	fmt.Println(cal)
}
